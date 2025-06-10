package model

type App struct {
	AppKey    string `gorm:"primarykey" json:"appKey"`
	AppSecret string `json:"appSecret"`

	ModelNoKey
}

func AddApp(app *App) error {
	return db.Create(app).Error
}

func UpdateApp(app *App) (bool, error) {
	r := db.Updates(app)
	return r.RowsAffected > 0, r.Error
}

func GetAppByAppKey(appkey string) (app App, err error) {
	err = db.Where(App{AppKey: appkey}).First(&app).Error
	return
}
