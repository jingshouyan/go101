package storage

import (
	"context"
	"fmt"
	"go101/config"
	"go101/model"
	"go101/response"
	"io"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	md5 := c.Query("md5")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	if md5 != "" {
		f, err := model.GetFileByMD5AndSize(md5, header.Size)
		if err == nil {
			f2 := &model.File{
				Name:       header.Filename,
				Extension:  getExtension(header.Filename),
				Size:       header.Size,
				MD5:        f.MD5,
				IsDir:      false,
				UploaderID: 0,
				Idx:        f.Idx,
			}
			model.AddFile(f2)
			response.OK(c, f2)
			return
		} else if err != gorm.ErrRecordNotFound {
			response.CommonError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	idx := uuid.NewString()
	ui, err := m.c.PutObject(context.Background(), m.cfg.Bucket, idx, file, header.Size, minio.PutObjectOptions{})
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	f2 := &model.File{
		Name:       header.Filename,
		Extension:  getExtension(header.Filename),
		Size:       ui.Size,
		MD5:        ui.ETag,
		IsDir:      false,
		UploaderID: 0,
		Idx:        idx,
	}
	model.AddFile(f2)
	response.OK(c, f2)
}

func (m *minioStorage) Download(c *gin.Context) {
	objectName := c.Query("file")
	if objectName == "" {
		response.CommonError(c, http.StatusBadRequest, "filename is empty")
		return
	}
	reader, err := m.c.GetObject(c, m.cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectName))
	c.Header("Content-Type", "application/octet-stream")

	// 将 MinIO 中的文件流式传输到客户端
	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}

}

func getExtension(fileName string) string {
	return path.Ext(fileName)
}
