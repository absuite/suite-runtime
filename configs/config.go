package configs

import (
	"os"

	"github.com/ggoop/goutils/glog"
	"github.com/ggoop/goutils/utils"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}
type DbConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ConnectionConfig struct {
	Driver   string `mapstructure:"driver"`
	Host     string `mapstructure:"host"`
	Port     uint   `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type Config struct {
	App AppConfig
	Db  DbConfig
}

const cmdRoot = "config"

var Default Config

func New() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(cmdRoot)
	viper.AddConfigPath(utils.JoinCurrentPath("env"))
	err := viper.ReadInConfig()
	if err != nil {
		glog.Errorf("Fatal error when reading %s config file:%s", cmdRoot, err)
		os.Exit(1)
	}

	viper.Unmarshal(&Default)
	if Default.App.Port == "" {
		Default.App.Port = "8080"
	}
}
