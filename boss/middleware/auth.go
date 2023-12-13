package middleware

import (
	"go101/model"
	"go101/response"
	"go101/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var permissions = util.NewSet[string]()

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		aid := getAdminIDFromContext(c)
		if aid == 0 {
			response.CommonError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()

			return
		}
		c.Keys["aid"] = aid
		c.Next()
	}
}

func Permission(requiredPermission string) gin.HandlerFunc {
	permissions.Add(requiredPermission)
	return func(c *gin.Context) {
		aid := c.Keys["aid"].(uint)
		admin, err := model.GetAdminById(aid)
		if err != nil || !admin.HasPermission(requiredPermission) {
			response.CommonError(c, http.StatusForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetAllPermissions() []string {
	return permissions.Values()
}

func getAdminIDFromContext(c *gin.Context) uint {
	session := sessions.Default(c)
	uid, _ := session.Get("aid").(uint)
	return uid
}
