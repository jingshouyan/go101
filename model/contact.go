package model

type Contact struct {
	Model
	UserID int64 `json:"userId"`

	ContactID int64  `json:"contactId"`
	Remark    string `json:"remark"`
}
