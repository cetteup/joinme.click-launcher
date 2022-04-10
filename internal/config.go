package internal

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

const (
	ConfigFilename = "config.yaml"
)

type Config struct {
	Games map[string]CustomLauncherConfig `yaml:"games"`
}

type CustomLauncherConfig struct {
	InstallPath    string   `yaml:"install_path"`
	ExecutablePath string   `yaml:"executable_path"`
	Args           []string `yaml:"args"`
}

func (c *Config) GetCustomLauncherConfig(game string) *CustomLauncherConfig {
	gameConfig, exists := c.Games[game]
	if !exists {
		return nil
	}
	return &gameConfig
}

func LoadConfig() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(wd, ConfigFilename)
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	err = yaml.Unmarshal(content, &RunningConfig)
	if err != nil {
		return err
	}

	return nil
}

var RunningConfig = &Config{
	Games: map[string]CustomLauncherConfig{},
}
