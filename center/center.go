package center

import (
	"context"
	"net/http"
	"sync"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/xgo11/env"
	"github.com/xgo11/stdlog"
)

import (
	"parallel_rsa/center/impl"
	"parallel_rsa/center/workerPool"
	"parallel_rsa/commonLib"
	gm "parallel_rsa/commonLib/ginMiddlewares"
)

type implCenter struct {
	sync.Once

	httpAddr string
	rpcAddr  string
	server   *http.Server
	mCh      chan struct{}
}

func NewServer(httpAddr string, rpcAddr string) commonLib.IServer {
	return &implCenter{httpAddr: httpAddr, rpcAddr: rpcAddr}
}

func (ic *implCenter) Run() {
	ic.Do(func() {
		if !env.ISDebug() {
			gin.SetMode(gin.ReleaseMode)
			gm.SetAvoidLogPath("/", "/favicon.ico")
		}
		var engine = gin.New()
		var wares []gin.HandlerFunc
		if env.ISDebug() {
			wares = []gin.HandlerFunc{gin.Logger(), gm.AppLogger, gin.Recovery()}
		} else {
			wares = []gin.HandlerFunc{gm.AppLogger, gm.SecurityInspector, gm.Recovery()}
		}
		engine.Use(wares...)
		impl.InitRoutes(engine)

		workerPool.InitGRPCAddr(ic.rpcAddr)

		ic.server = &http.Server{
			Addr:    ic.httpAddr,
			Handler: engine,
		}

		ic.mCh = make(chan struct{})

		go workerPool.Monitor(ic.mCh)

		stdlog.Std.Debugf("start running center http server ......")
		if err := ic.server.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				stdlog.Std.Panicf("center http serve stopped, %v", err)
			}
		}
	})

}

func (ic *implCenter) Shutdown() {
	close(ic.mCh)

	if ic.server != nil {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		if err := ic.server.Shutdown(ctx); err != nil {
			stdlog.Std.Errorf("center http server shutdown fail, %v", err)
		} else {
			stdlog.Std.Infof("center http server safe exit.")
		}
	}
}
