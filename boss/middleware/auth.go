package middleware

import (
	"go101/model"
	"go101/response"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement your authentication logic here
		// For example, check for a valid session or token
		uid := getAdminIDFromContext(c)
		if uid == 0 {
			response.CommonError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()

			return
		}
		c.Next()
	}
}

func Permission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getAdminIDFromContext(c)
		admin, err := model.GetAdminById(uid)
		if err != nil || !admin.HasPermission(requiredPermission) {
			response.CommonError(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}

func getAdminIDFromContext(c *gin.Context) uint {
	session := sessions.Default(c)
	uid, _ := session.Get("uid").(uint)
	return uid
}
