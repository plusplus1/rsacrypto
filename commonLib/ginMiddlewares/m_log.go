package ginMiddlewares

import (
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AddNoticeLog(ctx *gin.Context, k string, v interface{}) {
	var traceDict = ctx.GetStringMap(KeyTrace)
	if traceDict == nil {
		traceDict = make(map[string]interface{})
	}
	traceDict[k] = v
	ctx.Set(KeyTrace, traceDict)
}

func AppLogger(ctx *gin.Context) {
	var path = ctx.Request.URL.Path
	if _, ok := avoidPathMap[path]; ok { // skip log
		ctx.Next()
		return
	}

	// before
	_ = ctx.Request.ParseForm()

	startTime := time.Now()
	ctx.Set(keyStartTime, startTime)
	ctx.Set(KeyCode, 0)
	ctx.Set(KeyError, "")

	logFields := logrus.Fields{
		keyReqPath:   path,
		keyRemoteIP:  ctx.Request.RemoteAddr,
		keyReqMethod: ctx.Request.Method,
	}

	if ua := ctx.Request.UserAgent(); ua != "" {
		logFields[keyUserAgent] = ua
	}

	if referer := ctx.Request.Referer(); referer != "" {
		logFields[keyRefer] = referer
	}

	// params 数据太大，不记录日志
	//if argsLen := len(ctx.Request.Form); argsLen > 0 {
	//	params := make(map[string]interface{}, argsLen)
	//	for k, v := range ctx.Request.Form {
	//		if len(v) != 1 {
	//			params[k] = v
	//		} else {
	//			params[k] = v[0]
	//		}
	//	}
	//	argsBytes, _ := json.Marshal(params)
	//	logFields[keyParameters] = string(argsBytes)
	//}

	// process
	ctx.Next()

	// after
	logFields[KeyCode] = ctx.GetInt(KeyCode)
	if strError := ctx.GetString(KeyError); strError != "" {
		logFields[KeyError] = strError
	}
	for k, v := range ctx.GetStringMap(KeyTrace) {
		logFields[k] = v
	}

	logFields[keyCostMS] = fmt.Sprintf("%.3f", time.Since(startTime).Seconds()*1000)

	if statusCode := ctx.Writer.Status(); statusCode != 200 {
		logFields[keyStatus] = statusCode
	}
	logger.WithFields(logFields).Info()

}
