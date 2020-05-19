package httpLib

import (
	"fmt"
	"net/http"
)

const (
	Version     = "v0.1"
	ErrorOK     = http.StatusOK
	ErrorParams = http.StatusBadRequest // 参数不符合要求的都视为 bad request
	ErrorRTO    = http.StatusRequestTimeout
	ErrorServer = http.StatusInternalServerError
)

var (
	errTemplate = map[int]string{
		ErrorParams: "请求无效,%v",
		ErrorRTO:    "处理超时,%v",
		ErrorServer: "处理错误,%v",
	}
)

func FormatErrorMessage(code int, args ...interface{}) string {
	if tpl, ok := errTemplate[code]; ok {
		return fmt.Sprintf(tpl, args...)
	}
	return "UnknownError"
}
