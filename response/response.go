package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Success = 200

	ErrorExist    = 10001
	ErrorNotExist = 10002

	PasswordWrong = 20001
)

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": Success,
		"data": data,
	})
}

func BizError(c *gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
	})
}

func CommonError(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"code": code,
		"data": data,
	})
}
