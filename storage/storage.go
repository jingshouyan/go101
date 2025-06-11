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
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Storage interface {
	Save(fh *multipart.FileHeader, f *model.File) error
	Load(f *model.File) (io.ReadSeekCloser, int64, error)
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
	rsc, fileSize, err := storage.Load(&f)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	defer rsc.Close()

	encodedFileName := url.QueryEscape(f.Name)

	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// 设置响应头
		c.Header("Content-Length", fmt.Sprintf("%d", fileSize))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedFileName))
		c.Header("Content-Type", "application/octet-stream")
		c.Status(http.StatusOK)
		// 直接发送整个文件
		_, err = io.Copy(c.Writer, rsc)
		if err != nil {
			log.Warn("copy file error", zap.Error(err))
		}
		return
	}

	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
		response.CommonError(c, http.StatusRequestedRangeNotSatisfiable, err.Error())
		return
	}

	contentLength := end - start + 1
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", fmt.Sprintf("%d", contentLength))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedFileName))
	c.Status(http.StatusPartialContent)

	// 发送指定字节
	rsc.Seek(start, io.SeekStart)
	_, err = io.CopyN(c.Writer, rsc, contentLength)
	if err != nil {
		log.Warn("copy file error", zap.Error(err))
	}

}

// Range 格式示例：bytes=0-1023
func parseRange(header string, fileSize int64) (int64, int64, error) {
	if !strings.HasPrefix(header, "bytes=") {
		return 0, 0, fmt.Errorf("invalid range header")
	}

	rangeSpec := strings.TrimPrefix(header, "bytes=")
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	var start, end int64
	var err error

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start byte")
	}

	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid end byte")
		}
	}

	if start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("out of range")
	}

	return start, end, nil
}

func getExtension(fileName string) string {
	return path.Ext(fileName)
}
