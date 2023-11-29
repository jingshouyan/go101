package model

import (
	"go101/config"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"moul.io/zapgorm2"
)

var cfg = config.Conf.DB

var db *gorm.DB
var log = zap.L()

func init() {
	var err error
	gormLogger := zapgorm2.New(log)
	gormLogger.SetAsDefault()
	db, err = gorm.Open(sqlite.Open(cfg.Name), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
		Logger: gormLogger,
	})
	if err != nil {
		log.Error("failed to initialize database", zap.Error(err))
		panic(err)
	}

}
