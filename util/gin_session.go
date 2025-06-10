package util

import (
	"go101/g"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SaveAdminIDToSession(ctx *gin.Context, id int64) {
	session := sessions.Default(ctx)
	session.Set(g.AdminIdKey, id)
	session.Save()
}

func GetAdminIDFormSession(ctx *gin.Context) int64 {
	session := sessions.Default(ctx)
	uid, _ := session.Get(g.AdminIdKey).(int64)
	return uid
}

func ClearSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}
