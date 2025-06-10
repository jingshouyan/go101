package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	Logger  LoggerConfig  `json:"logger"`
	DB      DBConfig      `json:"db"`
	Storage StorageConfig `json:"storage"`
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
	PprofToken     string        `json:"pprof_token"`
}

type DBConfig struct {
	Driver      string `json:"driver"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	TablePrefix string `json:"table_prefix"`
}

const (
	StorageDriverMinio = "minio"
	StorageDriverLocal = "local"
)

type StorageConfig struct {
	Driver string             `json:"driver"`
	Minio  MinioConfig        `json:"minio"`
	Local  LocalStorageConfig `json:"local"`
}

type LocalStorageConfig struct {
	RootPath string `json:"root_path"`
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint"`
	UseSSL          bool   `json:"use_ssl"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
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
