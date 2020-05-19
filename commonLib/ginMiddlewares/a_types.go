package ginMiddlewares

import "parallel_rsa/commonLib/logLib"

const (
	KeyError = "error" // 错误信息
	KeyCode  = "code"  // 错误码
	KeyTrace = "trace" // 日志跟踪
)

const (
	keyParameters = "params"    // 请求参数
	keyStartTime  = "start"     // 开始时间
	keyCostMS     = "cost"      // 耗时/微秒
	keyUserAgent  = "ua"        // UserAgent
	keyRefer      = "refer"     // Referer
	keyReqPath    = "path"      // Request path
	keyRemoteIP   = "remote_ip" // Remote IP
	keyReqMethod  = "method"    // request method
	keyStatus     = "status"    // http status code
)

var (
	avoidPathMap = make(map[string]bool, 0)
	logger       = logLib.NewLogAdapter("HTTPSvr")
)

func SetAvoidLogPath(paths ...string) {
	for _, p := range paths {
		avoidPathMap[p] = true
	}
}
