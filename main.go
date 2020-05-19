package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"parallel_rsa/commonLib/logLib"
	"parallel_rsa/worker"
)

import (
	"github.com/plusplus1/logrus-extension"
	"github.com/xgo11/env"
	"github.com/xgo11/stdlog"
)

import (
	"parallel_rsa/center"
	"parallel_rsa/commonLib"
)

const (
	mCenter = "center"
	mWorker = "worker"
)

var (
	modeInitMap = map[string]func(){mWorker: initWorker, mCenter: initCenter}
	servers     = make([]commonLib.IServer, 0)
)

var (
	currentMode = ""
	httpSvr     = ""
	rpcSvr      = ""
)

func parseArgs() {
	flag.StringVar(&currentMode, "mode", "", "startup mode, e.g: "+strings.Join([]string{mCenter, mWorker}, "/"))
	flag.StringVar(&httpSvr, "http", "0.0.0.0:8000", "http server ")
	flag.StringVar(&rpcSvr, "rpc", "0.0.0.0:8001", "rpc server ")
	flag.Parse()

	if _, ok := modeInitMap[currentMode]; !ok {
		stdlog.Std.Errorf("mode invalid")
		flag.Usage()
		os.Exit(-1)
	}
	if httpSvr == "" {
		stdlog.Std.Errorf("http invalid")
		flag.Usage()
		os.Exit(-1)
	}
	if rpcSvr == "" {
		stdlog.Std.Errorf("rpc invalid")
		flag.Usage()
		os.Exit(-1)
	}
	logrus_extension.InitFileHook(filepath.Join(env.ConfDir(), fmt.Sprintf("logs/%s.yaml", currentMode)))
	stdlog.Std.Debugf("current mode=%v, httpSvr=%v, rpcSvr=%v", currentMode, httpSvr, rpcSvr)
}

func initCenter() {
	servers = append(servers, center.NewServer(httpSvr, rpcSvr))
}

func initWorker() {
	logger := logLib.NewLogAdapter("worker")
	servers = append(servers, worker.NewServer(httpSvr, rpcSvr, &logger))

}

func main() {
	parseArgs()
	if initFunc := modeInitMap[currentMode]; initFunc != nil {
		initFunc()
	}

	for _, s := range servers {
		go s.Run()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	for _, s := range servers {
		s.Shutdown()
	}
}
