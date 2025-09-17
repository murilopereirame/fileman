package config

import (
	"encoding/json"
	"fileman/fs"
)

type ConfigHandler struct {
	config string
}

func New(config string) *ConfigHandler {
	return &ConfigHandler{
		config: config,
	}
}

type IConfigHandler interface {
	Load() (Config, error)
}

func (h ConfigHandler) Load() (Config, error) {
	config := Config{}

	fileSystem := fs.FS{}
	content, readError := fileSystem.ReadFile(h.config)

	if readError != nil {
		return config, readError
	}

	err := json.Unmarshal(content, &config)

	return config, err
}

type WatchedDirectory struct {
	Path string
	Age  float64
}
type Config struct {
	Cron               string
	WatchedDirectories []WatchedDirectory
}
