package model

import "go101/util"

type PartFile struct {
	ModelStringKey
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MD5      string `json:"md5" gorm:"index"`
	PartSize int64  `json:"partSize"`
	Idx      string `json:"idx" gorm:"index"`
	Status   int64  `json:"status"` // 0: init, 1: uploading, 2: complete, 3: error
}

type PartFileChunk struct {
	ModelStringKey
	PartFileID string `json:"partFileId" gorm:"uniqueIndex:idx_part_file_chunk"`
	ChunkIndex int64  `json:"chunkIndex" gorm:"uniqueIndex:idx_part_file_chunk"`
	MD5        string `json:"md5"`
	Size       int64  `json:"size"`
	Idx        string `json:"idx"`
}

func AddPartFile(pf *PartFile) error {
	if pf.ID == "" {
		pf.ID = util.GenStringId()
	}
	return db.Create(pf).Error
}

func UpdatePartFile(pf *PartFile) error {
	return db.Save(pf).Error
}

func GetPartFileById(id string) (pf PartFile, err error) {
	err = db.First(&pf, map[string]interface{}{"id": id}).Error
	return
}

func AddPartFileChunk(pfc *PartFileChunk) error {
	if pfc.ID == "" {
		pfc.ID = util.GenStringId()
	}
	return db.Create(pfc).Error
}

func GetPartFileChunkByIndex(partFileID string, chunkIndex int64) (pfc PartFileChunk, err error) {
	err = db.First(&pfc, map[string]interface{}{"part_file_id": partFileID, "chunk_index": chunkIndex}).Error
	return
}

func GetPartFileChunksByPartFileID(partFileID string) (chunks []PartFileChunk, err error) {
	err = db.Where("part_file_id = ?", partFileID).Order("chunk_index asc").Find(&chunks).Error
	return
}

func UpdatePartFileChunk(pfc *PartFileChunk) error {
	return db.Save(pfc).Error
}
