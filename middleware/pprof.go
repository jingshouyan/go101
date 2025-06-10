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
		token := c.GetHeader(pprofTokenKey)
		if token == "" {
			token = c.Query(pprofTokenKey)
		}
		if token == config.Conf.Server.PprofToken {
			c.Next()
			return
		}
		response.CommonError(c, http.StatusForbidden)
		c.Abort()
	}
}
