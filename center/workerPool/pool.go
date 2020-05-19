package workerPool

import (
	"errors"
	"sync"
)

import (
	"parallel_rsa/commonLib/rpcLib"
)

type implPool struct{}

type Pool interface {
	Len() int
	Get() (rpcLib.RsaServiceClient, error)
	Release(rpcLib.RsaServiceClient, ...error)
}

var (
	inst         Pool = &implPool{} //ensure interface
	grpcAddrOnce      = &sync.Once{}
	grpcAddr          = ""
)

func PoolInst() Pool {
	return inst
}

func InitGRPCAddr(netAddr string) {
	if netAddr != "" {
		grpcAddrOnce.Do(func() {
			grpcAddr = netAddr
		})
	}
}

func ListEndpoints() []string {
	if grpcAddr != "" {
		return []string{grpcAddr}
	}
	return []string{}
}

func (p *implPool) Len() int {
	return 1
}

func (p *implPool) Get() (client rpcLib.RsaServiceClient, err error) {
	client, err = NewWkClient(grpcAddr)
	if client == nil && err == nil {
		err = errors.New("NoAvailableWorker")
	}
	return
}

func (p *implPool) Release(client rpcLib.RsaServiceClient, errs ...error) {
	if client == nil {
		return
	}

	var wk *wkClient
	if c, ok := client.(*wkClient); ok && c != nil {
		wk = c
	}

	if wk != nil && wk.conn != nil {
		_ = wk.conn.Close()
	}
}
