package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

import (
	"DCrypto/libs"
	"DCrypto/models"
)

func doDecrypt(c *gin.Context) {

	var err error
	var dataBytes []byte
	var rsaSK *rsa.PrivateKey

	req := ReqCryptOuter{}
	if err = c.Bind(&req); err != nil {
		c.String(400, "param error "+err.Error())
		return
	}
	if req.Data == "" || req.PriKey == "" {
		c.String(400, "param error")
		return
	}

	if rsaSK, err = libs.LoadPemPKCSPriKey(req.PriKey); err != nil {
		c.String(400, "PriKey error "+err.Error())
		return
	}

	if dataBytes, err = base64.StdEncoding.DecodeString(req.Data); err != nil {
		c.String(400, "Invalid base64 data error")
		return
	}

	var bitLen = rsaSK.N.BitLen()
	var dataLen = len(dataBytes)
	var step = 1024 * 128 // 128K 每个分片

	if step%bitLen != 0 {
		step = bitLen * (step / bitLen)
	}

	var totalCount = (dataLen + step - 1) / step
	var wks = manager.Choose(totalCount)
	if len(wks) < totalCount {
		c.String(500, "worker leak")
		return
	}

	var finishCount = 0
	var errChan = make(chan error)
	var doneChan = make(chan *models.Frame)

	bytesPriKey := []byte(req.PriKey)

	for i := 0; i < totalCount; i++ {

		var start = i * step
		var end = start + step
		if end > dataLen {
			end = dataLen
		}

		inputFrame := &models.Frame{
			Index:  i,
			Data:   dataBytes[start:end],
			PriKey: bytesPriKey,
		}

		go oneDecrypt(wks[i], inputFrame, doneChan, errChan)

	}

	var afterChan = time.After(15 * time.Second)

	var results = make([][]byte, totalCount)

	for {
		select {
		case f := <-doneChan:
			results[f.Index] = []byte(f.Data)
			finishCount++
		case e := <-errChan:
			c.String(500, e.Error())
			return
		case <-afterChan:
			c.String(500, "decrypt timeout")
			return
		}
		if finishCount >= totalCount {
			break
		}
	}

	for finishCount < totalCount {
		time.Sleep(10 * time.Microsecond)
	}

	c.String(200, string(bytes.Join(results, nil)))
	return
}

func oneDecrypt(w *models.Worker, f *models.Frame, doneChan chan *models.Frame, errChan chan error) {
	var err error
	var maxRetry = 3

	for i := 0; i < maxRetry; i++ {
		//fmt.Println("retry = ", i)
		err = nil
		if w == nil {
			if wks := manager.Choose(1); len(wks) < 1 {
				err = errors.New("worker leak")
				continue
			} else {
				w = wks[0]
			}
		}
		if bs, e := w.Decrypt(f, manager); e != nil {
			err = e
			logrus.Errorf("worker error, %v, error=%v", w.Unique(), e)
			w = nil
			continue
		} else {
			f.Data = bs
			doneChan <- f
			return
		}
	}
	if err != nil {
		errChan <- err
	}
}
