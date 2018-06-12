package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

import (
	"github.com/plusplus1/logrus-extension"
	"github.com/plusplus1/utils4g"
)

import (
	"DCrypto/models"
)

var (
	appConfig  *models.AppConf
	hostCenter string
)

func getLocalIp() string {
	arrAddress, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range arrAddress {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}

		}
	}
	return "localhost"
}

func reportHandle() func() {

	me := &models.Worker{IP: getLocalIp(), Port: appConfig.Worker.Port}
	log.Printf("[INFO]\treport worker =%s", me.Unique())

	addrCenter := appConfig.Worker.Center

	return func() {
		var conn net.Conn
		var err error

		if conn, err = net.Dial("tcp", addrCenter); err == nil {
			defer conn.Close()

			me.UpdateTime = time.Now().Unix()
			encoder := gob.NewEncoder(conn)
			if err = encoder.Encode(me); err == nil {
				_, err = conn.Read(make([]byte, 1))
			}
		}

		logrus.WithFields(logrus.Fields{
			"reqm":  "report",
			"error": err,
		}).Debug()
	}
}

func main() {

	mode := utils4g.DebugUtils.IsDebug()
	appConfig = &models.AppConf{}
	utils4g.ConfigUtils.ReadPath("app", appConfig)
	hostCenter = appConfig.Worker.Center

	log.Printf("[INFO]\tStart, debug mode = %v, conf = %v", mode, appConfig.Worker)
	logrus_extension.InitFileHook(filepath.Join(utils4g.ConfigUtils.GetBaseDir(), appConfig.Worker.Log))

	go func() {
		startTcpServer()
	}()

	tick := time.NewTicker(time.Duration(appConfig.Worker.Report) * time.Second)
	handler := reportHandle()
	handler()

	for {
		select {
		case <-tick.C:
			handler()
		}
	}

}
