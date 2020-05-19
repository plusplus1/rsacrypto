package httpLib

import (
	"github.com/gin-gonic/gin"
)

import (
	gm "parallel_rsa/commonLib/ginMiddlewares"
)

type Out struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

func NewOutputResult(code int, data interface{}, message ...interface{}) *Out {
	result := &Out{Code: code, Data: data}
	if result.Code != 0 && result.Code != ErrorOK {
		result.Message = FormatErrorMessage(code, message...)
	}
	return result
}

func (o *Out) beforeResp(ctx *gin.Context) {
	ctx.Set(gm.KeyCode, o.Code)
	if o.Code != 0 && o.Code != ErrorOK && o.Message != "" {
		ctx.Set(gm.KeyError, o.Message)
	}
}

func (o *Out) ResponseJSON(ctx *gin.Context) {
	o.beforeResp(ctx)
	ctx.JSON(ErrorOK, o)
}

func (o *Out) ResponseStatusWithMessage(ctx *gin.Context) {
	o.beforeResp(ctx)
	ctx.String(o.Code, o.Message)
}

func (o *Out) ResponseStatusWithDataString(ctx *gin.Context) {
	o.beforeResp(ctx)
	if s, ok := o.Data.(string); ok {
		ctx.String(ErrorOK, s)
		return
	}
	if bs, ok := o.Data.([]byte); ok {
		ctx.Status(ErrorOK)
		_, _ = ctx.Writer.Write(bs)
		ctx.Writer.Flush()
		return
	}
}
