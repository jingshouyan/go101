package main

import (
	"go101/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode("")
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
