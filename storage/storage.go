package storage

import (
	"crypto/md5"
	"fmt"
	"go101/config"
	"go101/model"
	"go101/response"
	"go101/util"
	"hash"
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

const (
	md5HeaderKey = "X-File-MD5"
)

type Storage interface {
	Save(fh *multipart.FileHeader, f *model.File) (string, error)
	Load(f *model.File) (io.ReadSeekCloser, int64, error)
	Delete(f *model.File) error
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
	if md5 == "" {
		md5 = c.GetHeader(md5HeaderKey)
	}
	if md5 == "" {
		md5 = c.PostForm("md5")
	}

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
	md5Server, err := storage.Save(fh, f2)
	if md5 != md5Server {
		log.Warn("md5 mismatch", zap.String("client", md5), zap.String("server", md5Server))
		f2.MD5 = md5Server // 更新实际存储的MD5
	}
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
		id = c.Param("id")
	}
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

	encodedfilename := url.QueryEscape(f.Name)

	rangeHeader := c.GetHeader("Range")
	// 未指定下载范围
	if rangeHeader == "" {
		// 设置响应头
		c.Header("Content-Length", fmt.Sprintf("%d", fileSize))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedfilename))
		c.Header("Content-Type", "application/octet-stream")
		c.Header(md5HeaderKey, f.MD5)
		c.Status(http.StatusOK)
		// 直接发送整个文件
		_, err = io.Copy(c.Writer, rsc)
		if err != nil {
			log.Error("copy file error", zap.Error(err))
		}
		return
	}

	// 指定下载范围
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
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedfilename))
	c.Header(md5HeaderKey, f.MD5)
	c.Status(http.StatusPartialContent)

	// 发送指定字节
	rsc.Seek(start, io.SeekStart)
	_, err = io.CopyN(c.Writer, rsc, contentLength)
	if err != nil {
		log.Error("copy file error", zap.Error(err))
	}

}

func Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.CommonError(c, http.StatusBadRequest, "id is empty")
		return
	}
	f, err := model.GetFileById(id)
	if err != nil {
		response.CommonError(c, http.StatusNotFound, err.Error())
		return
	}
	err = model.DeleteFile(id)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	count, err := model.CountFilesByIdx(f.Idx)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if count == 0 {
		// 当前文件是最后一个文件，删除文件
		err = storage.Delete(&f)
		if err != nil {
			log.Error("delete file error", zap.String("id", id), zap.Error(err))
		}
	}

	response.OK(c, nil)
}

type UploadInitReq struct {
	Filename string `json:"filename" binding:"required"`
	MD5      string `json:"md5" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
	PartSize int64  `json:"PartSize" binding:"required"`
}

type UploadInitRsp struct {
	State    int64          `json:"state"`
	File     model.File     `json:"file,omitempty"`
	PartFile model.PartFile `json:"partFile,omitempty"` // 如果是分片上传，返回分片文件信息
}

const (
	UploadStateInit      = 0
	UploadStateUploading = 1
	UploadStateComplete  = 2
	UploadStateError     = 3
)

func UploadInit(c *gin.Context) {
	// 处理上传初始化请求
	var req UploadInitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	// 检查文件是否已存在
	f, err := model.GetFileByMD5AndSize(req.MD5, req.Size)
	if err == nil {
		f2 := &model.File{
			Name:       req.Filename,
			Extension:  path.Ext(req.Filename),
			Size:       req.Size,
			MD5:        req.MD5,
			IsDir:      false,
			UploaderID: 0,
			Idx:        f.Idx,
		}
		model.AddFile(f2)
		rsp := UploadInitRsp{
			State: UploadStateComplete,
			File:  *f2,
		}
		response.OK(c, rsp)
		return
	}
	if err != gorm.ErrRecordNotFound {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	pf := &model.PartFile{
		Name:     req.Filename,
		Size:     req.Size,
		MD5:      req.MD5,
		PartSize: req.PartSize,
		Idx:      util.GenStringId(),
		Status:   UploadStateInit,
	}

	// TODO: storage 初始化分片

	err = model.AddPartFile(pf)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	rsp := UploadInitRsp{
		State:    UploadStateUploading,
		PartFile: *pf,
	}
	response.OK(c, rsp)

}

type UploadPartReq struct {
	PartFileID string `json:"partFileId" binding:"required"`
	ChunkIndex int64  `json:"chunkIndex" binding:"required"`
	MD5        string `json:"md5" binding:"required"`
}

func UploadPart(c *gin.Context) {
	var req UploadPartReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	p, err := model.GetPartFileChunkByIndex(req.PartFileID, req.ChunkIndex)
	if err != nil && err != gorm.ErrRecordNotFound {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err == nil && (p.MD5 == req.MD5 && p.Size == fh.Size) {
		// 分片已存在，直接返回
		response.OK(c, nil)
		return
	}
	// TODO： storage 保存分片，校验md5

	if err == gorm.ErrRecordNotFound {
		// 未找到记录
		p2 := &model.PartFileChunk{
			PartFileID: req.PartFileID,
			ChunkIndex: req.ChunkIndex,
			Size:       fh.Size,
			MD5:        req.MD5,
			Idx:        util.GenStringId(),
		}
		err = model.AddPartFileChunk(p2)
		if err != nil {
			response.CommonError(c, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// 更新已存在的分片
		p.Size = fh.Size
		p.MD5 = req.MD5
		p.Idx = util.GenStringId()
		err = model.UpdatePartFileChunk(&p)
		if err != nil {
			response.CommonError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.OK(c, nil)

}

type UploadCompleteReq struct {
	PartFileID string `json:"partFileId" binding:"required"`
}

func UploadComplete(c *gin.Context) {
	var req UploadCompleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CommonError(c, http.StatusBadRequest, err.Error())
		return
	}
	pf, err := model.GetPartFileById(req.PartFileID)
	if err == gorm.ErrRecordNotFound {
		response.CommonError(c, http.StatusNotFound, "part file not found")
		return
	}
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if pf.Status == UploadStateComplete {
		// 分片上传已完成，直接返回文件信息
		f, err := model.GetFileById(pf.ID)
		if err != nil {
			response.CommonError(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.OK(c, f)
		return
	}
	pfcs, err := model.GetPartFileChunksByPartFileID(req.PartFileID)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(pfcs) != int(pf.PartSize) {
		response.BizError(c, response.UploadUncomplete, "not all parts uploaded")
		return
	}
	// TODO: storage 合并分片
	pf.Status = UploadStateComplete
	model.UpdatePartFile(&pf)
	f := &model.File{
		ModelStringKey: model.ModelStringKey{
			ID: pf.ID,
		},
		Name:       pf.Name,
		Extension:  getExtension(pf.Name),
		Size:       pf.Size,
		MD5:        pf.MD5,
		IsDir:      false,
		UploaderID: 0,
		Idx:        pf.Idx,
	}
	err = model.AddFile(f)
	if err != nil {
		response.CommonError(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, f)
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

func getExtension(filename string) string {
	return path.Ext(filename)
}

func newHasher() hash.Hash {
	return md5.New()
}
