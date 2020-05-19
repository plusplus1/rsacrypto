package worker

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/xgo11/env"
	"github.com/xgo11/stdlog"
	"google.golang.org/grpc"
)

import (
	"parallel_rsa/commonLib"
	"parallel_rsa/commonLib/logLib"
	"parallel_rsa/commonLib/rpcLib"
	"parallel_rsa/worker/impl_v1"
)

type implWorker struct {
	sync.Once

	httpAddr   string
	rpcAddr    string
	logger     *logLib.Adapter
	httpServer *http.Server
	grpcServer *grpc.Server
}

func NewServer(httpSvr string, rpcSvr string, logger *logLib.Adapter) commonLib.IServer {
	return &implWorker{httpAddr: httpSvr, rpcAddr: rpcSvr, logger: logger}
}

func (iw *implWorker) Run() {
	iw.Do(func() {
		if !env.ISDebug() {
			gin.SetMode(gin.ReleaseMode)
		}
		var engine = gin.Default()

		engine.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "pong")
		})

		engine.GET("/list", func(ctx *gin.Context) {
			if iw.grpcServer == nil {
				ctx.String(http.StatusOK, "grpc server has not started!")
				return
			}
			var methods []string

			for _, info := range iw.grpcServer.GetServiceInfo() {
				for _, m := range info.Methods {
					methods = append(methods, m.Name)
				}
			}
			var out = struct {
				Methods []string `json:"methods"`
				Count   int      `json:"count"`
			}{Methods: methods, Count: len(methods)}

			ctx.JSON(http.StatusOK, &out)
			return
		})

		iw.httpServer = &http.Server{Addr: iw.httpAddr, Handler: engine}

		iw.grpcServer = grpc.NewServer()
		// register rpc methods
		rpcLib.RegisterRsaServiceServer(iw.grpcServer, impl_v1.NewServer(iw.logger))

		go func() {
			stdlog.Std.Debugf("start running worker http server ......")
			if err := iw.httpServer.ListenAndServe(); err != nil {
				if err.Error() != "http: Server closed" {
					stdlog.Std.Panicf("worker http serve stopped, %v", err)
				}
			}
		}()

		go func() {
			if lis, err := net.Listen("tcp", iw.rpcAddr); err == nil {
				stdlog.Std.Infof("grpc serve %s , start ......", iw.rpcAddr)
				if err = iw.grpcServer.Serve(lis); err != nil {
					log.Panicf("grpc serve at %s failed, error=%v", iw.rpcAddr, err)
				}
			} else {
				stdlog.Std.Panicf("start listen at %s failed, error=%v", iw.rpcAddr, err)
			}
		}()
	})

}

func (iw *implWorker) Shutdown() {
	if iw.httpServer != nil {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		if err := iw.httpServer.Shutdown(ctx); err == nil {
			stdlog.Std.Infof("worker http server safe stopped!")
		} else {
			stdlog.Std.Panicf("worker http server shutdown fail, %v!", err)
		}
	}

	if iw.grpcServer != nil {
		iw.grpcServer.GracefulStop()
		stdlog.Std.Infof("worker grpc server graceful stopped!")
	}
}
