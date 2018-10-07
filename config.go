package main

import (
	"log"
	"os"
	"strings"

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

const cmdRoot = "config"

var appConfig AppConfig
var dbConfig DbConfig

func init() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName(cmdRoot)
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error when reading %s config file:%s", cmdRoot, err)
		os.Exit(1)
	}

	viper.UnmarshalKey("app", &appConfig)
	if appConfig.Port == "" {
		appConfig.Port = "8080"
	}
	viper.UnmarshalKey("db", &dbConfig)
}
