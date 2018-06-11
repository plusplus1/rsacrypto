package RSAUtils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func DecodePem(key string) (*pem.Block, error) {
	return DecodePemBytes([]byte(key))
}

func DecodePemBytes(key []byte) (block *pem.Block, err error) {
	if p, _ := pem.Decode(key); p != nil {
		return p, nil
	} else {
		return nil, errors.New("decode fail")
	}
}

func DecodePublicKeyBytes(key []byte) (pk *rsa.PublicKey, err error) {

	var block *pem.Block
	if block, err = DecodePemBytes(key); err == nil {
		var k interface{}
		if k, err = x509.ParsePKIXPublicKey(block.Bytes); err == nil {
			if pk, ok := k.(*rsa.PublicKey); ok {
				return pk, nil
			}
		}
	}
	return
}

func DecodePKCSPrivateKeyBytes(key []byte) (sk *rsa.PrivateKey, err error) {
	var block *pem.Block
	if block, err = DecodePemBytes(key); err != nil {
		return
	}

	if sk, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		var tmp interface{}
		if tmp, err = x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			sk = tmp.(*rsa.PrivateKey)
		}
	}
	if err != nil {
		return nil, err
	}
	if sk == nil {
		err = errors.New("parse fail")
	}
	return
}

func DecodePublicKey(key string) (pk *rsa.PublicKey, err error) {
	return DecodePublicKeyBytes([]byte(key))
}
func DecodePKCSPrivateKey(key string) (sk *rsa.PrivateKey, err error) {
	return DecodePKCSPrivateKeyBytes([]byte(key))
}
