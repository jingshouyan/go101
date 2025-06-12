package storage

import (
	"context"
	"encoding/hex"
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

func (s *minioStorage) Save(fh *multipart.FileHeader, f *model.File) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := newHasher()
	_, err = s.c.PutObject(context.Background(), s.cfg.Bucket, f.Idx, io.TeeReader(file, hasher), fh.Size, minio.PutObjectOptions{})
	hashSum := hex.EncodeToString(hasher.Sum(nil))
	return hashSum, err
}

func (s *minioStorage) Load(f *model.File) (io.ReadSeekCloser, int64, error) {
	rsc, err := s.c.GetObject(context.Background(), s.cfg.Bucket, f.Idx, minio.GetObjectOptions{})
	if err != nil {
		return nil, 0, err
	}
	info, err := rsc.Stat()
	return rsc, info.Size, err
}

func (s *minioStorage) Delete(f *model.File) error {
	err := s.c.RemoveObject(context.Background(), s.cfg.Bucket, f.Idx, minio.RemoveObjectOptions{})
	if err != nil {
		log.Error("minio delete object error", zap.String("bucket", s.cfg.Bucket), zap.String("object", f.Idx), zap.Error(err))
		return err
	}
	log.Info("minio object deleted", zap.String("bucket", s.cfg.Bucket), zap.String("object", f.Idx))
	return nil
}
