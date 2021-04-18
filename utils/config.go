package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	IsProduction bool
	InputFile    string
	UseBin       bool
	UseLog       bool
	PrintInput   bool
	MasterId     int32
	MasterPort   int32
	ObserverId   int32
	ObserverPort int32
}

var LoadedConfig Config

func ReloadConfig(filePath string) (error, Config) {
	if filePath == "" {
		filePath = "config.json"
	}
	file, _ := os.Open(filePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
		return err, config
	}

	LoadedConfig = config

	return nil, config
}
