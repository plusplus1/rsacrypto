package rsaLib

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"io"
	"runtime"
)

const (
	deParallelSize = 1024 * 20
)

type onePart struct {
	Index int
	Data  []byte
}

type sha1HashOpts struct{}

func (so sha1HashOpts) HashFunc() crypto.Hash {
	return crypto.SHA1
}

func CalcDecryptGroupSize(sk *rsa.PrivateKey) int {
	return sk.N.BitLen() / 8
}

func CalcEncryptGroupSize(pk *rsa.PublicKey) int {
	return pk.N.BitLen()/8 - 28
}

func PureEncrypt(data []byte, pk *rsa.PublicKey) ([]byte, error) {

	bufInput := bytes.NewBuffer(data)
	bufOutput := bytes.NewBuffer(nil)

	groupSize := CalcEncryptGroupSize(pk)
	bufGroup := make([]byte, groupSize)

	for {
		if n, errReadBuf := bufInput.Read(bufGroup); errReadBuf != nil && errReadBuf != io.EOF {
			return nil, errReadBuf
		} else if n > 0 {
			if bs, errEncrypt := rsa.EncryptPKCS1v15(rand.Reader, pk, bufGroup[0:n]); errEncrypt != nil {
				return nil, errEncrypt
			} else {
				if _, errWriteBuf := bufOutput.Write(bs); errWriteBuf != nil {
					return nil, errWriteBuf
				}
			}
		} else {
			break
		}
	}
	return bufOutput.Bytes(), nil

}

func PureDecrypt(cipher []byte, sk *rsa.PrivateKey, runParallel ...bool) ([]byte, error) {
	bufInput := bytes.NewBuffer(cipher)
	var shouldParallel = false
	if len(runParallel) >= 1 {
		shouldParallel = runParallel[0]
	} else if len(cipher) > deParallelSize {
		shouldParallel = true
	}

	if shouldParallel {
		return decryptParallel(bufInput, sk)
	}
	return decryptOneProcess(bufInput, sk)
}

func decryptOneProcess(bufInput *bytes.Buffer, sk *rsa.PrivateKey) ([]byte, error) {

	var groupSize = CalcDecryptGroupSize(sk)
	var bufGroup = make([]byte, groupSize)
	var bufOutput = bytes.NewBuffer(nil)

	for {
		if n, errReadBuf := bufInput.Read(bufGroup); errReadBuf != nil && errReadBuf != io.EOF {
			return nil, errReadBuf
		} else if n > 0 {
			if bs, errDecrypt := rsa.DecryptPKCS1v15(rand.Reader, sk, bufGroup[0:n]); errDecrypt != nil {
				return nil, errDecrypt
			} else {
				if _, errWriteBuf := bufOutput.Write(bs); errWriteBuf != nil {
					return nil, errWriteBuf
				}
			}
		} else {
			break
		}
	}
	return bufOutput.Bytes(), nil
}

func decryptParallel(bufInput *bytes.Buffer, sk *rsa.PrivateKey) ([]byte, error) {

	inSize := bufInput.Len()

	var groupSize = CalcDecryptGroupSize(sk)
	var totalCount = (inSize + groupSize - 1) / groupSize
	var finishedCount = 0

	var taskChan = make(chan *onePart, runtime.NumCPU()-1)
	var doneChan = make(chan error, 1)

	var outParts = make([][]byte, totalCount)

	decryptOnePart := func(in *onePart) {
		bs, err := rsa.DecryptPKCS1v15(rand.Reader, sk, in.Data)
		if err == nil {
			outParts[in.Index] = bs
		}
		doneChan <- err
	}

	go func() {
		for {
			t := <-taskChan
			go decryptOnePart(t)
		}

	}()

	go func() {
		for i := 0; i < totalCount; i++ {
			var buf = make([]byte, groupSize)
			n, _ := bufInput.Read(buf)
			one := &onePart{Index: i, Data: buf[0:n]}
			taskChan <- one
		}
	}()

	for {
		select {

		case err := <-doneChan:
			finishedCount++
			if err != nil {
				return nil, err
			}
		}

		if finishedCount >= totalCount {
			break
		}
	}

	return bytes.Join(outParts, nil), nil
}

func Sign(text string, sk *rsa.PrivateKey) (string, error) {

	h := sha1HashOpts{}.HashFunc()
	hObj := h.New()
	hObj.Write([]byte(text))
	hashed := hObj.Sum(nil)

	if bs, err := rsa.SignPKCS1v15(rand.Reader, sk, h, hashed); err != nil {
		return "", err
	} else {
		return base64.StdEncoding.EncodeToString(bs), nil
	}

}

func Verify(text string, sig string, pk *rsa.PublicKey) bool {
	if sigBytes, err := base64.StdEncoding.DecodeString(sig); err == nil {
		h := sha1HashOpts{}.HashFunc()
		hObj := h.New()
		hObj.Write([]byte(text))
		hashed := hObj.Sum(nil)

		if err = rsa.VerifyPKCS1v15(pk, h, hashed, sigBytes); err == nil {
			return true
		}
	}
	return false
}
