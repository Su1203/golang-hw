package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf
	HTTP     HTTPConf
	Database DatabaseConf
	Storage  StorageConf
}

type LoggerConf struct {
	Level string
}

type HTTPConf struct {
	Host string
	Port string
}

type DatabaseConf struct {
	DSN string
}

type StorageConf struct {
	Type string
}

func NewConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	ext := strings.ToLower(filepath.Ext(configPath))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".toml":
		if err := toml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse TOML config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format: %s (supported: .yaml, .yml, .toml)", ext)
	}

	if config.Logger.Level == "" {
		config.Logger.Level = "INFO"
	}
	if config.HTTP.Host == "" {
		config.HTTP.Host = "localhost"
	}
	if config.HTTP.Port == "" {
		config.HTTP.Port = "8080"
	}
	if config.Storage.Type == "" {
		config.Storage.Type = "memory"
	}

	return &config, nil
}
