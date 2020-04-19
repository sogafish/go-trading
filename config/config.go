package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey    string
	ApiSecret string
	LogFile   string
}

var Config ConfigList

func init() {
	_config, err := ini.Load("config.ini")

	if err != nil {
		log.Printf("Fail To Read Config File: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		ApiKey:    _config.Section("bitflyer").Key("api_key").String(),
		ApiSecret: _config.Section("bitflyer").Key("api_secret").String(),
		LogFile:   _config.Section("go-trading").Key("log_file").String(),
	}
}
