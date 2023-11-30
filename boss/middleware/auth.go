package middleware

import (
	"go101/model"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implement your authentication logic here
		// For example, check for a valid session or token
		uid := getAdminIDFromContext(c)
		if uid == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}
		c.Next()
	}
}

func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getAdminIDFromContext(c)
		admin, err := model.GetAdminById(uid)
		if err != nil || !admin.HasPermission(requiredPermission) {
			c.AbortWithStatus(http.StatusForbidden)
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
