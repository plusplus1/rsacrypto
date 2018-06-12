package main

import (
	"github.com/gin-gonic/gin"
)

import (
	"DCrypto/middlewares"
)

func initHttpServer() (engine *gin.Engine) {

	engine = gin.New()
	wares := []gin.HandlerFunc{
		gin.Logger(),                  // 请求日志
		middlewares.SecurityInspector, // 流量安全检查
		middlewares.AppLogger,         // 应用日志
		gin.Recovery(),                // recovery
	}

	engine.Use(wares...)

	pingHandler := func(c *gin.Context) {
		c.Writer.WriteString(APP_VERSION)
	}
	engine.GET("/dcrypto/ping", pingHandler)
	engine.GET("/dcrypto/workers/list", listWorkers)
	engine.POST("/dcrypto/decrypt", doDecrypt)
	return
}
