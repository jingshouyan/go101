package router

import (
	"go101/boss/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var p = middleware.Permission

func InitRoute(r *gin.RouterGroup) {
	store := cookie.NewStore([]byte("adwoeriumliuu"))
	r.Use(sessions.Sessions("boss", store))
	r.POST("/login", login)
	r.GET("/logout", logout)
	r.Use(middleware.Auth())
	r.POST("/changePwd", changePwd)
	r.GET("/profile", getProfile)
	r.PUT("/profile", updateProfile)
	r.GET("/ping", p("ping"), ping)

	r.GET("/admin", p("admin:list"), pageQueryAdmin)
	r.GET("/admin/:id", p("admin:get"), getAdminById)
	r.POST("/admin", p("admin:add"), addAdmin)
	r.PUT("/admin/:id", p("admin:edit"), editAdmin)
	r.DELETE("/admin/:id", p("admin:delete"), deleteAdmin)
	r.PUT("/admin/:id/state", p("admin:state"), changeAdminState)

}
