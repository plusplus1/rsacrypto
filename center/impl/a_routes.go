package impl

import (
	"github.com/gin-gonic/gin"
)

import (
	"parallel_rsa/commonLib/httpLib"
)

func InitRoutes(engine *gin.Engine) {
	engine.GET("/", httpLib.Ping)
	engine.GET("/workers/list", doListWorkers)
	engine.POST("/rsa/decrypt", doRsaDecrypt)
	engine.POST("/rsa/encrypt", doRsaEncrypt)
}
