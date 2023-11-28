package router

import (
	"go101/config"
	"go101/logger"

	"github.com/gin-gonic/gin"
)

var cfg = config.Conf.Server

func Serve() {
	gin.SetMode(cfg.Mode)
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(cfg.Addr)
}
