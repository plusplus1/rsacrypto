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
	"DCrypto/models"
)

func startTcpServer() {

	var tcpAddr *net.TCPAddr
	var port = appConfig.Center.TcpPort
	var strAddress = fmt.Sprintf("0.0.0.0:%v", port)

	tcpAddr, _ = net.ResolveTCPAddr("tcp", strAddress)
	if tcpListener, err := net.ListenTCP("tcp", tcpAddr); err == nil {
		log.Printf("[INFO]\tStart tcp server at %v success!", strAddress)
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

func handle(conn *net.TCPConn) {

	if conn == nil {
		return
	}

	worker := models.Worker{}
	decoder := gob.NewDecoder(conn)
	var err error

	if err = decoder.Decode(&worker); err == nil {
		worker.UpdateTime = time.Now().Unix()
		manager.Add(&worker)
		conn.Write([]byte("ok"))
	} else {
		conn.Write([]byte(err.Error()))
	}

	logrus.WithFields(logrus.Fields{
		"reqm":   "report",
		"worker": worker.Unique(),
		"error":  err,
	}).Debug()

}
