package impl_v1

import (
	"context"
	"strings"
	"time"
)

import (
	"github.com/sirupsen/logrus"
)

import (
	"parallel_rsa/commonLib/logLib"
	"parallel_rsa/commonLib/rpcLib"
	"parallel_rsa/commonLib/rsaLib"
)

type serverV1 struct {
	logger *logLib.Adapter
}

func NewServer(logger ...*logLib.Adapter) *serverV1 {
	if len(logger) > 0 && logger[0] != nil {
		return &serverV1{logger: logger[0]}
	}
	lg := logLib.NewLogAdapter("rpc")
	return &serverV1{logger: &lg}
}

func (s *serverV1) startTrace() logrus.Fields {
	trace := logrus.Fields{"_st": time.Now()}
	return trace
}

func (s *serverV1) endTrace(trace logrus.Fields, args ...interface{}) {
	nt := make(logrus.Fields)
	if v, ok := trace["_st"]; ok {
		if st, ok := v.(time.Time); ok {
			nt["cost"] = time.Now().Sub(st).Truncate(time.Millisecond)
		}
	}

	for k, v := range trace {
		if strings.HasPrefix(k, "_") {
			continue
		}
		nt[k] = v
	}

	for _, arg := range args {
		if arg == nil {
			continue
		}
		if resp, ok := arg.(*rpcLib.RsaResponse); ok && resp != nil {
			code, msg := resp.GetCode(), resp.GetMsg()
			nt["code"] = code
			nt["seq"] = resp.GetSeqNo()
			if code != 0 || msg != "" {
				nt["error"] = msg
			}
		}
	}

	s.logger.WithFields(nt).Info()
}

func (s *serverV1) DoEncrypt(ctx context.Context, req *rpcLib.RsaRequest) (resp *rpcLib.RsaResponse, err error) {

	var trace = s.startTrace()

	trace["method"] = "DoEncrypt"
	trace["seq"] = req.GetSeqNo()

	if pubKey, pkErr := rsaLib.LoadPemPubKeyBytes(req.GetKey()); pkErr != nil {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Code: 101, Msg: pkErr.Error()}
	} else if result, encErr := rsaLib.PureEncrypt(req.GetBody(), pubKey); encErr != nil {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Code: 110, Msg: encErr.Error()}
	} else {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Data: result}
	}

	s.endTrace(trace, resp)
	return
}

func (s *serverV1) DoDecrypt(ctx context.Context, req *rpcLib.RsaRequest) (resp *rpcLib.RsaResponse, err error) {

	var trace = s.startTrace()

	trace["method"] = "DoDecrypt"
	trace["seq"] = req.GetSeqNo()

	if priKey, pkErr := rsaLib.LoadPemPKCSPriKeyBytes(req.GetKey()); pkErr != nil {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Code: 101, Msg: pkErr.Error()}
	} else if result, decErr := rsaLib.PureDecrypt(req.GetBody(), priKey); decErr != nil {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Code: 110, Msg: decErr.Error()}
	} else {
		resp = &rpcLib.RsaResponse{SeqNo: req.GetSeqNo(), Data: result}
	}

	s.endTrace(trace, resp)
	return
}

func (s *serverV1) Ping(ctx context.Context, req *rpcLib.HbRequest) (resp *rpcLib.HbResponse, err error) {
	resp = new(rpcLib.HbResponse)
	resp.Msg = "pong"
	return
}
