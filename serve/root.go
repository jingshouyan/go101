package serve

import (
	bossrouter "go101/boss/router"
	"go101/config"
	"go101/middleware"
	"go101/storage"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var cfg = config.Conf.Server
var log = zap.L()

func Serve() {

	gin.SetMode(cfg.Mode)
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	addRouter(r)

	s := &http.Server{
		Addr:           cfg.Addr,
		Handler:        r,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
	log.Info("server start", zap.String("addr", cfg.Addr))
	err := s.ListenAndServe()
	log.Info("server stop", zap.Error(err))
}

func addRouter(r *gin.Engine) {
	dg := r.Group("/debug")
	dg.Use(middleware.PprofAuth())
	pprof.RouteRegister(dg, "/pprof")
	bossGroup := r.Group("/boss")
	bossrouter.InitRoute(bossGroup)

	fileGroup := r.Group("/file")
	storage.InitRoute(fileGroup)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
