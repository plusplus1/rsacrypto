package workerPool

import (
	"context"
)
import (
	"google.golang.org/grpc"
)
import (
	"parallel_rsa/commonLib/rpcLib"
)

type wkClient struct {
	connStr string
	conn    *grpc.ClientConn
	client  rpcLib.RsaServiceClient
}

func NewWkClient(connStr string) (*wkClient, error) {
	if conn, connErr := grpc.Dial(connStr, grpc.WithInsecure()); connErr != nil {
		return nil, connErr
	} else {
		return &wkClient{connStr: connStr, conn: conn}, nil
	}
}

func (c *wkClient) ensureClient() {
	if c.client == nil {
		c.client = rpcLib.NewRsaServiceClient(c.conn)
	}
}

// rsa加密
func (c *wkClient) DoEncrypt(ctx context.Context, in *rpcLib.RsaRequest, opts ...grpc.CallOption) (*rpcLib.RsaResponse, error) {
	c.ensureClient()
	return c.client.DoEncrypt(ctx, in, opts...)
}

// rsa解密
func (c *wkClient) DoDecrypt(ctx context.Context, in *rpcLib.RsaRequest, opts ...grpc.CallOption) (*rpcLib.RsaResponse, error) {
	c.ensureClient()
	return c.client.DoDecrypt(ctx, in, opts...)
}

// 心跳检测服务
func (c *wkClient) Ping(ctx context.Context, in *rpcLib.HbRequest, opts ...grpc.CallOption) (*rpcLib.HbResponse, error) {
	c.ensureClient()
	return c.Ping(ctx, in, opts...)
}
