package model

type Role struct {
	Model

	RoleName string      `json:"roleName"`
	Codes    StringSlice `json:"codes"`
	All      bool        `json:"all"`
}

func AddRole(r *Role) error {
	return db.Create(r).Error
}

func UpdateRole(r *Role) (bool, error) {
	r1 := db.Updates(r)
	return r1.RowsAffected > 0, r1.Error
}

func DeleteRole(id uint) error {
	return db.Delete(&Role{}, id).Error
}

func GetRoleById(id uint) (r Role, err error) {
	err = db.First(&r, id).Error
	return
}

func GetRoles(limit, offset int, maps interface{}) (roles []Role, err error) {
	err = db.Limit(limit).Offset(offset).Find(&roles).Error
	return
}

func CountRoles(maps interface{}) (count int64, err error) {
	err = db.Model(&Role{}).Where(maps).Count(&count).Error
	return
}
