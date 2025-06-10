package storage

import (
	"go101/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Storage interface {
	Upload(c *gin.Context)
	Download(c *gin.Context)
}

var log = zap.L()

func NewStorage() Storage {
	driver := config.Conf.Storage.Driver
	switch driver {
	case "minio":
		return NewMinioStorage()
	case "local":
		return NewLocalStorage()
	}

	log.Panic("storage driver not found", zap.String("driver", driver))
	return nil
}
