package storage

import (
	"context"
	"go101/config"
	"go101/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type minioStorage struct {
	c   *minio.Client
	cfg config.MinioConfig
}

func NewMinioStorage() *minioStorage {
	cfg := config.Conf.Storage.Minio
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	}

	c, err := minio.New(cfg.Endpoint, options)
	if err != nil {
		log.Panic("minio init error", zap.Error(err))
	}
	exist, err := c.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		log.Panic("minio bucket exists error", zap.Error(err))
	}
	if !exist {
		err = c.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{Region: cfg.Region})
		if err != nil {
			log.Panic("minio make bucket error", zap.Error(err))
		}
		log.Info("minio bucket created", zap.String("bucket", cfg.Bucket))
	} else {
		log.Info("minio bucket exists", zap.String("bucket", cfg.Bucket))
	}
	log.Info("minio init success")

	return &minioStorage{c: c, cfg: cfg}
}

func (m *minioStorage) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()
	_, err = m.c.PutObject(context.Background(), m.cfg.Bucket, header.Filename, file, header.Size, minio.PutObjectOptions{})
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}
