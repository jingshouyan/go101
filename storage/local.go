package storage

import (
	"go101/config"

	"github.com/gin-gonic/gin"
)

type LocalStorage struct {
	RootPath string
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		RootPath: config.Conf.Storage.Local.RootPath,
	}
}

func (s *LocalStorage) Upload(c *gin.Context) {
}

func (s *LocalStorage) Download(c *gin.Context) {
}
