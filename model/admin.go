package model

import (
	"errors"

	"gorm.io/gorm/clause"
)

const (
	AccountTypePhone int = iota + 1
	AccountTypeEmail
	AccountTypeUsername
)

type Admin struct {
	Model

	Phone      string `json:"phone" gorm:"unique"`
	Email      string `json:"email" gorm:"unique"`
	Username   string `json:"username" gorm:"unique"`
	PwdHash    string `json:"-"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	RoleID     uint   `json:"roleId"`
	Role       Role   `json:"role"`
	DisabledAt int64  `json:"disabledAt" gorm:"default:-1"`
}

func (a *Admin) HasPermission(requiredPermission string) bool {
	if a.DisabledAt > 0 {
		return false
	}
	if a.Role.All {
		return true
	}
	for _, code := range a.Role.Codes {
		if code == requiredPermission {
			return true
		}
	}
	return false
}

func AddAdmin(a *Admin) error {
	return db.Create(a).Error
}

func UpdateAdmin(a *Admin) (bool, error) {
	r := db.Updates(a)
	return r.RowsAffected > 0, r.Error
}

func DeleteAdmin(id uint) error {
	return db.Delete(&Admin{}, id).Error
}

func GetAdminById(id uint) (a Admin, err error) {
	err = db.First(&a, id).Error
	return
}

func GetAdmins(limit, offset int, maps interface{}) (admins []Admin, err error) {
	err = db.Preload(clause.Associations).Limit(limit).Offset(offset).Find(&admins).Error
	return
}

func CountAdmins(maps interface{}) (count int64, err error) {
	err = db.Model(&Admin{}).Where(maps).Count(&count).Error
	return
}

func GetAdminByAccount(t int, account string) (a Admin, err error) {
	c := new(Admin)
	switch t {
	case AccountTypePhone:
		c.Phone = account
	case AccountTypeEmail:
		c.Email = account
	case AccountTypeUsername:
		c.Username = account
	default:
		return a, errors.New("invalid account type")
	}
	err = db.Preload(clause.Associations).Where(c).First(&a).Error
	return
}
