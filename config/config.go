package config

import (
	"encoding/json"
	"fmt"

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
	Mode string `json:"mode"`
	Addr string `json:"addr"`
}

var Conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&Conf)
	c, _ := json.Marshal(Conf)
	fmt.Println(string(c))
}
