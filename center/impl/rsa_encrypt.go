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

func doRsaEncrypt(ctx *gin.Context) {

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
	if params.PubKey == "" {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "pub_key empty")
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var pubKey *rsa.PublicKey
	if pubKey, err = rsaLib.LoadPemPubKey(params.PubKey); err != nil {
		out := httpLib.NewOutputResult(httpLib.ErrorParams, nil, "pub_key invalid")
		out.ResponseStatusWithMessage(ctx)
		return
	}

	var dataBytes = []byte(params.Data)
	var datLen = len(dataBytes)

	var groupSize = 1024 * 256
	if gL := pubKey.N.BitLen() - 28; groupSize%gL != 0 {
		groupSize = gL * (groupSize / gL)
	}
	var groupCount = (datLen + groupSize - 1) / groupSize // 计算分组数量
	gm.AddNoticeLog(ctx, "batch", groupCount)

	// encrypt one part
	// TODO 重试优化, grpcClient 内部默认重试一次
	var encryptOnePart = func(ctxObj context.Context, req *rpcLib.RsaRequest, rCH chan *rpcLib.RsaResponse, eCH chan error) {
		var client rpcLib.RsaServiceClient
		var resp *rpcLib.RsaResponse
		var err error

		if client, err = workerPool.PoolInst().Get(); err == nil {
			defer workerPool.PoolInst().Release(client)
			resp, err = client.DoEncrypt(ctxObj, req)
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

	// loop encrypt for each part
	for i := 0; i < groupCount; i++ {
		start, end := i*groupSize, (i+1)*groupSize
		if end >= datLen {
			end = datLen
		}

		req := new(rpcLib.RsaRequest)
		req.SeqNo = int32(i)
		req.Key = []byte(params.PubKey)
		req.Body = dataBytes[start:end]

		go encryptOnePart(contextObj, req, resCH, errCH)
	}

	var resultErrCode = 0
	var resultErr error
	var resultBytesArr = make([][]byte, groupCount, groupCount)
	var finishCount = 0

	for resultErr == nil && finishCount < groupCount {
		select {

		case <-contextObj.Done():
			resultErrCode = httpLib.ErrorRTO
			resultErr = ctx.Err()
			cancelFunc()
		case err := <-errCH:
			resultErrCode = httpLib.ErrorServer
			resultErr = err

		case resp := <-resCH:
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
		rasEncryptedBytes := bytes.Join(resultBytesArr, nil)
		b64EncodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(rasEncryptedBytes)))
		base64.StdEncoding.Encode(b64EncodedBytes, rasEncryptedBytes)

		out := httpLib.NewOutputResult(0, b64EncodedBytes)
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
