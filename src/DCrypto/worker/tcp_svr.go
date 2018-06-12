package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)
import (
	"github.com/sirupsen/logrus"
)

import (
	"DCrypto/libs"
	"DCrypto/models"
)

func startTcpServer() {

	var tcpAddr *net.TCPAddr
	var port = appConfig.Worker.Port
	var strAddress = fmt.Sprintf("0.0.0.0:%v", port)

	tcpAddr, _ = net.ResolveTCPAddr("tcp", strAddress)
	if tcpListener, err := net.ListenTCP("tcp", tcpAddr); err == nil {
		log.Printf("[INTO]\tTcp server ready , at %v ...", strAddress)
		defer tcpListener.Close()

		for {
			tcpConn, err := tcpListener.AcceptTCP()
			if err != nil {
				continue
			}

			go handle(tcpConn)
		}
	} else {
		log.Panic(err)
	}
}

func handle(tcpConn *net.TCPConn) {
	if tcpConn == nil {
		return
	}

	var st = time.Now()
	var inputFrame = models.Frame{}
	var resultFrame = models.FrameRet{}

	var okChan = make(chan int)
	var timeAfterChan = time.After(5 * time.Second)

	go func(inf *models.Frame, out *models.FrameRet) {
		var decoder = gob.NewDecoder(tcpConn)
		if err := decoder.Decode(inf); err == nil {
			out.Index = inf.Index
			if sk, err := libs.LoadPemPKCSPriKeyBytes(inf.PriKey); err == nil {
				if result, err := libs.DecryptNoneBase64(inf.Data, sk, true); err == nil {
					out.Data = result
				} else {
					out.Error = err.Error()
				}
			} else {
				out.Error = err.Error()
			}

		} else {
			out.Error = err.Error()
		}
		okChan <- 1
	}(&inputFrame, &resultFrame)

	var hasTimeout bool
	select {
	case <-okChan:
		hasTimeout = false
	case <-timeAfterChan:
		hasTimeout = true
	}

	if hasTimeout {
		resultFrame.Error = "timeout"
	}

	encoder := gob.NewEncoder(tcpConn)
	encoder.Encode(&resultFrame)

	logrus.WithFields(logrus.Fields{
		"data_len": len(inputFrame.Data),
		"key_len":  len(inputFrame.PriKey),
		"index":    inputFrame.Index,
		"cost":     time.Since(st).Seconds() * 1000,
		"error":    resultFrame.Error,
	}).Info()

	return
}
