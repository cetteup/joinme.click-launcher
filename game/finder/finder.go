package finder

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
)

type Type string

const (
	RegistryFinder   Type = "RegistryFinder"
	CustomPathFinder Type = "CustomPathFinder"
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type Config struct {
	ForType           Type
	RegistryPath      string
	RegistryValueName string
	CustomInstallPath string
}

type GameFinder struct {
	repository RegistryRepository
}

func NewGameFinder(repository RegistryRepository) *GameFinder {
	return &GameFinder{
		repository: repository,
	}
}

func (f GameFinder) IsGameInstalledAnywhere(configs []Config) (bool, error) {
	for i, config := range configs {
		installed, err := f.IsGameInstalled(config)
		// Fail silently unless config is the last chance to find the game
		if err != nil && i == len(configs)-1 {
			return false, err
		}

		if installed {
			return true, nil
		}
	}

	return false, nil
}

func (f GameFinder) IsGameInstalled(config Config) (bool, error) {
	switch config.ForType {
	case CustomPathFinder:
		return f.isGameInstalledAccordingToCustomPath(config)
	default:
		return f.isGameInstalledAccordingToRegistry(config)
	}
}

func (f GameFinder) isGameInstalledAccordingToRegistry(config Config) (bool, error) {
	_, err := f.getGameInstallDirFromRegistry(config)

	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f GameFinder) isGameInstalledAccordingToCustomPath(config Config) (bool, error) {
	_, err := os.Stat(config.CustomInstallPath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f GameFinder) GetGameInstallDirFromSomewhere(configs []Config) (string, error) {
	for i, config := range configs {
		installDir, err := f.GetGameInstallDir(config)
		// Fail silently unless config is the last chance to find the game
		if err != nil && i == len(configs)-1 {
			return "", err
		}

		if installDir != "" {
			return installDir, nil
		}
	}

	return "", fmt.Errorf("game install path could not be determined")
}

func (f GameFinder) GetGameInstallDir(config Config) (string, error) {
	switch config.ForType {
	case CustomPathFinder:
		return config.CustomInstallPath, nil
	default:
		return f.getGameInstallDirFromRegistry(config)
	}
}

func (f GameFinder) getGameInstallDirFromRegistry(config Config) (string, error) {
	return f.repository.GetStringValue(registry.LOCAL_MACHINE, config.RegistryPath, config.RegistryValueName)
}
