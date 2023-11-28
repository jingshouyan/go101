package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `json:"server"`
	Logger LoggerConfig `json:"logger"`
}

type LoggerConfig struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"max_size"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
	Compress   bool   `json:"compress"`
	Level      string `json:"level"`
}

type ServerConfig struct {
	Mode           string        `json:"mode"`
	Addr           string        `json:"addr"`
	ReadTimeout    time.Duration `json:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes"`
}

var Conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		panic(err)
	}
}
