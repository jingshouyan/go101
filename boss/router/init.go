package router

import (
	"go101/boss/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var log = zap.L()

func InitRoute(r *gin.RouterGroup) {
	store := cookie.NewStore([]byte("adwoeriumliuu"))
	r.Use(sessions.Sessions("boss", store))
	r.POST("/login", login)
	r.GET("/logout", logout)
	r.Use(middleware.AuthMiddleware())
	r.GET("/ping", ping)
}
