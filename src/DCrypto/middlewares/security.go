package middlewares

import (
	"github.com/gin-gonic/gin"
)

// SecurityInspector, 流量基础安全检查
func SecurityInspector(c *gin.Context) {

	//c.Request.Host
	c.Next()

}
