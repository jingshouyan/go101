package middleware

import (
	"go101/g"
	"go101/model"
	"go101/response"
	"go101/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

var permissions = util.NewSet[string]()

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		aid := util.GetAdminIDFormSession(c)
		if aid == 0 {
			response.CommonError(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()

			return
		}
		c.Set(g.AdminIdKey, aid)
		c.Set(g.AuthType, g.AuthTypeSession)
		c.Next()
	}
}

func Permission(requiredPermission string) gin.HandlerFunc {
	permissions.Add(requiredPermission)
	return func(c *gin.Context) {
		aid := util.GetAdminIDFormSession(c)
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
