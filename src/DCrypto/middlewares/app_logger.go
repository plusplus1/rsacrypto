package middlewares

import (
	"fmt"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	KEY_PATH       = "url"       // 请求url
	KEY_START_TIME = "start"     // 开始时间
	KEY_END_TIME   = "end"       // 结束处理时间
	KEY_COST       = "cost"      // 耗时/微秒
	KEY_UA         = "ua"        // UserAgent
	KEY_REFER      = "refer"     // Referer
	KEY_REMOTE_IP  = "remote_ip" // Remote IP
	KEY_REQ_METHOD = "reqm"      // request method

	// 返回信息
	KEY_ERRORNO   = "error"     // 错误信息
	KEY_MSG       = "code"      // 错误码
	KEY_ADDITIONS = "additions" // 业务自定义扩展信息
)

func AppLogger(c *gin.Context) {

	// before
	c.Request.ParseForm()

	tStart := time.Now()
	c.Set(KEY_START_TIME, tStart)
	c.Set(KEY_ADDITIONS, make(map[string]interface{}))
	c.Set(KEY_MSG, 200)
	c.Set(KEY_ERRORNO, "")

	c.ClientIP()
	appInfo := logrus.Fields{
		KEY_REFER:      c.Request.Referer(),
		KEY_UA:         c.Request.UserAgent(),
		KEY_REMOTE_IP:  c.Request.RemoteAddr,
		KEY_REQ_METHOD: c.Request.Method,
		KEY_PATH:       c.Request.URL.Path,
	}

	// process
	c.Next()

	// after

	appInfo[KEY_MSG] = c.GetInt(KEY_MSG)
	if strError := c.GetString(KEY_ERRORNO); strError != "" {
		appInfo[KEY_ERRORNO] = strError
	}

	if additions := c.GetStringMap(KEY_ADDITIONS); len(additions) > 0 {
		for k, v := range additions {
			appInfo[k] = v
		}
	}

	appInfo[KEY_COST] = fmt.Sprintf("%.3f", time.Since(tStart).Seconds()*1000)
	logrus.WithFields(appInfo).Info()

}
