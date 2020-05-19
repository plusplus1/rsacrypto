package workerPool

import (
	"context"
	"time"
)
import (
	"github.com/xgo11/stdlog"
	"google.golang.org/grpc"
)
import (
	"parallel_rsa/commonLib/rpcLib"
)

func pingWorker() bool {

	if conn, connErr := grpc.Dial(grpcAddr, grpc.WithInsecure()); connErr == nil {
		defer func() {
			_ = conn.Close()
		}()

		client := rpcLib.NewRsaServiceClient(conn)
		if resp, reqErr := client.Ping(context.Background(), new(rpcLib.HbRequest)); reqErr == nil {
			if resp != nil && resp.Msg != "" {
				stdlog.Std.Infof("check worker heartbeat ok")
				return true
			}
		}
	}
	stdlog.Std.Errorf("check worker heartbeat fail")
	return false
}

func Monitor(ctrl chan struct{}) {

	var lstCheck = time.Now().Add(-60 * time.Minute)

	var intervalOk = 30 * time.Second
	var intervalFail = 3 * time.Second
	var loopSleep = 100 * time.Millisecond

	var interval time.Duration = 0

	for {
		select {
		case <-ctrl:
			return
		default:
			if time.Now().Sub(lstCheck) > interval {
				if ok := pingWorker(); ok {
					interval = intervalOk
				} else {
					interval = intervalFail
				}
				lstCheck = time.Now()
			}
		}
		time.Sleep(loopSleep)

	}

}
