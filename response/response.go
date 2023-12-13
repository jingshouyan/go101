package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Success = 200

	ErrorExist    = 10001
	ErrorNotExist = 10002

	PasswordWrong   = 20001
	AccountDisabled = 20002
)

func OK(c *gin.Context, data ...interface{}) {
	if len(data) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": Success,
			"data": data[0],
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": Success,
		})
	}

}

func BizError(c *gin.Context, code int, data ...interface{}) {
	if len(data) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"data": data[0],
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": code,
		})
	}
}

func CommonError(c *gin.Context, code int, data ...interface{}) {
	if len(data) > 0 {
		c.JSON(code, gin.H{
			"code": code,
			"data": data,
		})
	} else {
		c.JSON(code, gin.H{
			"code": code,
		})
	}
}
