package middleware

import (
	"go101/config"
	"go101/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	pprofTokenKey = "Pprof-Token"
)

func PprofAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader(pprofTokenKey) != config.Conf.Server.PprofToken {
			response.CommonError(c, http.StatusForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
