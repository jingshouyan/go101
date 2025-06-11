package storage

import (
	"context"
	"go101/config"
	"go101/model"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type minioStorage struct {
	c   *minio.Client
	cfg config.MinioConfig
}

func newMinioStorage() *minioStorage {
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

func (s *minioStorage) Save(fh *multipart.FileHeader, f *model.File) error {
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = s.c.PutObject(context.Background(), s.cfg.Bucket, f.Idx, file, fh.Size, minio.PutObjectOptions{})
	return err
}

func (s *minioStorage) Load(f *model.File) (io.ReadSeekCloser, int64, error) {
	rsc, err := s.c.GetObject(context.Background(), s.cfg.Bucket, f.Idx, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}
	info, err := rsc.Stat()
	return rsc, info.Size, err
}
