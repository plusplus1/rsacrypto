package ginMiddlewares

import (
	"github.com/gin-gonic/gin"
)

type recoveryWriter struct{}

func (rw *recoveryWriter) Write(p []byte) (n int, err error) {
	logger.Errorf(string(p))
	return 0, nil
}

func Recovery() func(c *gin.Context) {
	return gin.RecoveryWithWriter(new(recoveryWriter))
}
