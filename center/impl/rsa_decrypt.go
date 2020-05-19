package impl

import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)
import (
	"parallel_rsa/center/workerPool"
	gm "parallel_rsa/commonLib/ginMiddlewares"
	"parallel_rsa/commonLib/httpLib"
	"parallel_rsa/commonLib/rpcLib"
	"parallel_rsa/commonLib/rsaLib"
)

func doRsaDecrypt(ctx *gin.Context) {

	var params = rsaParams{}
	var err error

	if err = ctx.Bind(&params); err != nil {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "参数解析错误:"+err.Error())
		out.ResponseStatusWithMessage(ctx)
		return
	}

	if params.Data == "" {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "data empty")
		out.ResponseStatusWithMessage(ctx)
		return

	}
	if params.PriKey == "" {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "pri_key empty")
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var priKey *rsa.PrivateKey
	if priKey, err = rsaLib.LoadPemPKCSPriKey(params.PriKey); err != nil {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "pri_key invalid")
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var b64DecryptedDataBytes []byte
	if b64DecryptedDataBytes, err = base64.StdEncoding.DecodeString(params.Data); err != nil {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "data b64 decode fail")
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var datLen = len(b64DecryptedDataBytes)
	var groupSize = 1024 * 128                      // 128K 每个分片
	if kL := priKey.N.BitLen(); groupSize%kL != 0 { // 确保batchSize 是 keyLength的整数倍
		groupSize = kL * (groupSize / kL)
	}
	var groupCount = (datLen + groupSize - 1) / groupSize // 计算分组数量
	gm.AddNoticeLog(ctx, "batch", groupCount)

	// 解密一个分组
	// TODO: 优化重试，grpcClient 内部默认回重试一次
	var decryptOnePart = func(ctxObj context.Context, req *rpcLib.RsaRequest, rCH chan *rpcLib.RsaResponse, eCH chan error) {
		var client rpcLib.RsaServiceClient
		var err error
		var resp *rpcLib.RsaResponse

		if client, err = workerPool.PoolInst().Get(); err == nil {
			defer workerPool.PoolInst().Release(client)
			resp, err = client.DoDecrypt(ctxObj, req)
		}
		if resp != nil && err == nil {
			rCH <- resp
		} else {
			eCH <- err
		}
	}

	resCH := make(chan *rpcLib.RsaResponse, groupCount)
	errCH := make(chan error)

	contextObj, cancelFunc := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancelFunc()

	// loop request decrypt for each group
	for i := 0; i < groupCount; i++ {
		start, end := i*groupSize, (i+1)*groupSize
		if end >= datLen {
			end = datLen
		}

		req := new(rpcLib.RsaRequest)
		req.SeqNo = int32(i)
		req.Key = []byte(params.PriKey)
		req.Body = b64DecryptedDataBytes[start:end]

		go decryptOnePart(contextObj, req, resCH, errCH)
	}

	var resultErrCode = 0
	var resultErr error
	var resultBytesArr = make([][]byte, groupCount, groupCount)
	var finishCount = 0

	// wait for results
	for resultErr == nil && finishCount < groupCount {
		select {
		case <-contextObj.Done(): // when timeout , return err
			resultErrCode = httpLib.ErrorRTO
			resultErr = contextObj.Err()
			cancelFunc()
		case err := <-errCH: // when receive err, return
			resultErrCode = httpLib.ErrorServer
			resultErr = err
		case resp := <-resCH: // when receive response obj
			if resp.Code == 0 {
				resultBytesArr[int(resp.GetSeqNo())] = resp.GetData()
				finishCount++
			} else { // 处理错误的情况
				resultErrCode = httpLib.ErrorServer
				resultErr = fmt.Errorf("%v|%v|%v", resp.SeqNo, resp.Code, resp.Msg)
			}
		}
	}

	if finishCount == groupCount { // success
		out := httpLib.NewOutputResult(0, bytes.Join(resultBytesArr, nil))
		out.ResponseStatusWithDataString(ctx)
		return
	}

	if resultErr != nil {
		out := httpLib.NewOutputResult(resultErrCode, nil, resultErr)
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var out = httpLib.NewOutputResult(httpLib.ErrorServer, nil, "")
	out.ResponseStatusWithMessage(ctx)
	return
}
