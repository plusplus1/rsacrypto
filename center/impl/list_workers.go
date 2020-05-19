package impl

import (
	"github.com/gin-gonic/gin"
)

import (
	"parallel_rsa/center/workerPool"
	"parallel_rsa/commonLib/httpLib"
)

func doListWorkers(context *gin.Context) {
	out := httpLib.NewOutputResult(0, workerPool.ListEndpoints())
	out.ResponseJSON(context)

}
