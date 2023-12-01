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
	r.POST("/profile", updateProfile)
	r.GET("/ping", p("ping"), ping)

}
