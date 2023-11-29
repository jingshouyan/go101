package model

import (
	"go101/config"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

var cfg = config.Conf.DB

var db *gorm.DB
var log = zap.L()

func init() {
	var err error
	glog := zapgorm2.Logger{
		ZapLogger:                 log,
		LogLevel:                  logger.Warn,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          true,
		IgnoreRecordNotFoundError: true,
		Context:                   nil,
	}
	glog.SetAsDefault()
	db, err = gorm.Open(sqlite.Open(cfg.Name), &gorm.Config{
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
}

func migrate() {
	db.AutoMigrate(&Admin{})
	db.AutoMigrate(&Role{})
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
