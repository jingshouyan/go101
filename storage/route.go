package storage

import "github.com/gin-gonic/gin"

func InitRoute(r *gin.RouterGroup) {
	r.POST("/upload", Upload)
	r.GET("/download", Download)
	r.GET("/download/:id", Download)
	r.DELETE("/delete/:id", Delete)
}
