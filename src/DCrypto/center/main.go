package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

import (
	"github.com/plusplus1/logrus-extension"
	"github.com/plusplus1/utils4g"
)

import (
	"DCrypto/models"
)

const (
	APP_VERSION = "v1.0"
)

var (
	appConfig *models.AppConf
)

func main() {

	mode := utils4g.DebugUtils.IsDebug()
	appConfig = &models.AppConf{}
	utils4g.ConfigUtils.ReadPath("app", appConfig)

	logrus_extension.InitFileHook(filepath.Join(utils4g.ConfigUtils.GetBaseDir(), appConfig.Center.Log))

	log.Printf("[INFO]\tStart, debug mode = %v, conf = %v\n", mode, appConfig.Center)

	go func() {
		svr := initHttpServer()
		if err := svr.Run(fmt.Sprintf(":%v", appConfig.Center.HttpPort)); err != nil {
			log.Panic(err)
		}
	}()

	go func() {
		startTcpServer()
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Kill, os.Interrupt)

	<-sigChan
	shutdown()

}

func shutdown() {
	log.Print("[INFO]\texiting ...")
}
