package model

import (
	"database/sql"
	"fmt"
	"go101/config"
	"go101/util"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"

	_ "modernc.org/sqlite"
)

var cfg = config.Conf.DB

var db *gorm.DB
var log = zap.L()

func init() {
	var err error
	glog := zapgorm2.Logger{
		ZapLogger:                 log,
		LogLevel:                  logger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: true,
		Context:                   nil,
	}
	glog.SetAsDefault()
	sqlDB, err := sql.Open("sqlite", cfg.Name)
	if err != nil {
		log.Fatal("failed to open DB:", zap.Error(err))
	}

	db, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
		Logger: glog,
	})
	if err != nil {
		log.Error("failed to initialize database", zap.Error(err))
		panic(err)
	}
	migrate()
	initRole()
	initAdmin()
}

func migrate() {
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&File{})

	db.AutoMigrate(&App{})
	db.AutoMigrate(&User{})
}

func initRole() {
	r, err := GetRoleById(1)
	if err == gorm.ErrRecordNotFound {
		r = Role{
			Model:    Model{ID: 1},
			RoleName: "admin",
			Codes:    []string{"a", "b"},
			All:      true,
		}
		AddRole(&r)
	}
}

func initAdmin() {
	r, err := GetAdminById(1)
	if err == gorm.ErrRecordNotFound {
		usename := "admin"
		pwd, _ := util.GenerateRandomPassword(12)
		pwdHash, _ := util.HashPassword(pwd)
		r = Admin{
			Model:    Model{ID: 1},
			Username: usename,
			PwdHash:  pwdHash,
			RoleID:   1,
		}
		AddAdmin(&r)
		fmt.Printf("init admin{username:%s,password:%s},you shoud change it\n", usename, pwd)
		log.Info("init admin", zap.String("username", usename), zap.String("password", pwd))
	}
}
