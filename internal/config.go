package internal

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ConfigFilename = "config.yaml"
)

type config struct {
	DebugLogging bool                            `yaml:"debug_logging"`
	QuietLaunch  bool                            `yaml:"quiet_launch"`
	Games        map[string]CustomLauncherConfig `yaml:"games"`
}

func (c *config) GetCustomLauncherConfig(game string) *CustomLauncherConfig {
	gameConfig, exists := c.Games[game]
	if !exists {
		return nil
	}
	return &gameConfig
}

type CustomLauncherConfig struct {
	ExecutableName string   `yaml:"executable_name"`
	ExecutablePath string   `yaml:"executable_path"`
	InstallPath    string   `yaml:"install_path"`
	Args           []string `yaml:"args"`
}

func (c *CustomLauncherConfig) HasValues() bool {
	return c != nil && (c.HasExecutableName() || c.HasExecutablePath() || c.HasInstallPath() || c.HasArgs())
}

func (c *CustomLauncherConfig) HasExecutableName() bool {
	return c != nil && c.ExecutableName != ""
}

func (c *CustomLauncherConfig) HasExecutablePath() bool {
	return c != nil && c.ExecutablePath != ""
}

func (c *CustomLauncherConfig) HasInstallPath() bool {
	return c != nil && c.InstallPath != ""
}

func (c *CustomLauncherConfig) HasArgs() bool {
	return c != nil && len(c.Args) > 0
}

func LoadConfig() error {
	wd, err := os.Executable()
	if err != nil {
		return err
	}

	path := filepath.Join(filepath.Dir(wd), ConfigFilename)
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	err = yaml.Unmarshal(content, &Config)
	if err != nil {
		return err
	}

	return nil
}

var Config = &config{}
