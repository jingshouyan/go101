package storage

import (
	"fmt"
	"go101/config"
	"go101/model"
	"go101/response"
	"go101/util"
	"io"
	"mime/multipart"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage interface {
	Save(fh *multipart.FileHeader, f *model.File) error
	Load(f *model.File) (io.ReadCloser, error)
}

var log = zap.L()

var storage Storage = newStorage()

func newStorage() Storage {
	driver := config.Conf.Storage.Driver
	switch driver {
	case "minio":
		return newMinioStorage()
	case "local":
		return newLocalStorage()
	}

	log.Panic("storage driver not found", zap.String("driver", driver))
	return nil
}

func Upload(c *gin.Context) {
	md5 := c.Query("md5")
	fh, err := c.FormFile("file")
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	if md5 != "" {
		f, err := model.GetFileByMD5AndSize(md5, fh.Size)
		if err == nil {
			f2 := &model.File{
				Name:       fh.Filename,
				Extension:  getExtension(fh.Filename),
				Size:       fh.Size,
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

	idx := util.GenStringId()

	f2 := &model.File{
		Name:       fh.Filename,
		Extension:  getExtension(fh.Filename),
		Size:       fh.Size,
		MD5:        md5,
		IsDir:      false,
		UploaderID: 0,
		Idx:        idx,
	}
	err = storage.Save(fh, f2)
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	model.AddFile(f2)
	response.OK(c, f2)
}

func Download(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.CommonError(c, http.StatusBadRequest, "filename is empty")
		return
	}
	f, err := model.GetFileById(id)
	if err != nil {
		response.CommonError(c, http.StatusNotFound, err.Error())
		return
	}
	reader, err := storage.Load(&f)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer reader.Close()

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Name))
	c.Header("Content-Type", "application/octet-stream")

	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}

}

func getExtension(fileName string) string {
	return path.Ext(fileName)
}
