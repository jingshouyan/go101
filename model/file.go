package model

import (
	"go101/util"
)

type File struct {
	ModelStringKey
	Name       string `json:"name"`
	Extension  string `json:"extension"`
	Size       int64  `json:"size"`
	MD5        string `json:"md5" gorm:"index"`
	IsDir      bool   `json:"isDir"`
	UploaderID int64  `json:"uploaderId"`
	Idx        string `json:"idx" gorm:"index"`
}

func AddFile(f *File) error {
	f.ID = util.GenStringId()
	return db.Create(f).Error
}

func GetFileById(id string) (f File, err error) {
	err = db.First(&f, map[string]interface{}{"id": id}).Error
	return
}

func GetFileByMD5AndSize(md5 string, size int64) (f File, err error) {
	err = db.Where("md5 = ? AND size = ?", md5, size).First(&f).Error
	return
}

func DeleteFile(id string) error {
	return db.Where("id = ?", id).Delete(&File{}).Error
}

func CountFilesByIdx(idx string) (count int64, err error) {
	err = db.Model(&File{}).Where("idx = ?", idx).Count(&count).Error
	return
}
