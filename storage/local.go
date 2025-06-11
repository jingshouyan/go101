package storage

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"go101/config"
	"go101/model"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/howeyc/crc16"
	"go.uber.org/zap"
)

type LocalStorage struct {
	RootPath string
}

func newLocalStorage() *LocalStorage {
	path := config.Conf.Storage.Local.RootPath
	if path == "" {
		path = "./uploads"
	}
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Panic("create local storage root path error", zap.Error(err))
	}

	return &LocalStorage{
		RootPath: path,
	}
}

func (s *LocalStorage) Save(fh *multipart.FileHeader, f *model.File) (string, error) {
	pd := s.getParentDir(f.Idx)
	if err := os.MkdirAll(pd, 0755); err != nil {
		log.Error("create parent directory error", zap.String("path", pd), zap.Error(err))
		return "", err
	}
	filePath := fmt.Sprintf("%s/%s", pd, f.Idx)
	file, err := os.Create(filePath)
	if err != nil {
		log.Error("create file error", zap.String("path", filePath), zap.Error(err))
		return "", err
	}
	defer file.Close()
	f2, err := fh.Open()
	if err != nil {
		log.Error("open file header error", zap.String("filename", fh.Filename), zap.Error(err))
		return "", err
	}
	defer f2.Close()
	hasher := newHasher()
	_, err = io.Copy(file, io.TeeReader(f2, hasher))
	if err != nil {
		log.Error("copy file error", zap.String("from", fh.Filename), zap.String("to", filePath), zap.Error(err))
		return "", err
	}
	hashSum := hex.EncodeToString(hasher.Sum(nil))
	return hashSum, nil
}

func (s *LocalStorage) Load(f *model.File) (io.ReadSeekCloser, int64, error) {
	pd := s.getParentDir(f.Idx)
	fp := filepath.Join(pd, f.Idx)
	file, err := os.Open(fp)
	if err != nil {
		return nil, 0, err
	}
	info, err := file.Stat()
	return file, info.Size(), err
}

func (s *LocalStorage) getParentDir(filename string) string {
	crc := crc16.Checksum([]byte(filename), crc16.IBMTable)
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, crc)

	return fmt.Sprintf("%s/%02x/%02x/", s.RootPath, bytes[0], bytes[1])

}
