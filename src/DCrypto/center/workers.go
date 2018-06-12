package main

import (
	"encoding/json"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"DCrypto/models"
)

var (
	manager = models.NewWorkerManager()
)

func listWorkers(c *gin.Context) {
	manager.RLock()
	defer manager.RUnlock()

	bytes, _ := json.Marshal(manager.GetAllWorkers())
	c.Writer.Write(bytes)
}
