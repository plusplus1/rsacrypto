package httpLib

import (
	"github.com/gin-gonic/gin"
)

func Ping(context *gin.Context) {
	context.String(ErrorOK, Version)
}
