package model

type User struct {
	Model
	AppKey  string `json:"appKey" grom:"uniqueIndex:idx_app_account"`
	Account string `json:"account" grom:"uniqueIndex:idx_app_account"`
	Token   string `json:"token"`

	UserInfo
	UserConfig
}

type UserInfo struct {
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Sign      string `json:"sign"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Gender    int    `json:"gender"`
	Extension string `json:"extension"`
}

type UserConfig struct {
	Enabled bool `json:"enabled"`
}

func AddUser(u *User) error {
	return db.Create(u).Error
}

func UpdateUser(u *User) (bool, error) {
	r := db.Updates(u)
	return r.RowsAffected > 0, r.Error
}

func DeleteUserByID(id uint) error {
	return db.Delete(&User{}, id).Error
}

func DeleteUserByAppKeyAccount(appkey, account string) error {
	return db.Where(User{AppKey: appkey, Account: account}).Delete(&User{}).Error
}

func GetUserByID(id uint) (u User, err error) {
	err = db.First(&u, id).Error
	return
}

func GetUserByAppKeyAccount(appkey, account string) (u User, err error) {
	err = db.Where(User{AppKey: appkey, Account: account}).First(&u).Error
	return
}

func GetUsers(limit, offset int, maps interface{}) (users []User, err error) {
	err = db.Limit(limit).Offset(offset).Where(maps).Find(&users).Error
	return
}

func CountUsers(maps interface{}) (count int64, err error) {
	err = db.Model(&User{}).Where(maps).Count(&count).Error
	return
}
