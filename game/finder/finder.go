package finder

import (
	"errors"
	"fmt"
	"github.com/cetteup/joinme.click-launcher/internal"
	"golang.org/x/sys/windows/registry"
	"path/filepath"
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

type SoftwareFinder struct {
	repository RegistryRepository
}

func NewSoftwareFinder(repository RegistryRepository) *SoftwareFinder {
	return &SoftwareFinder{
		repository: repository,
	}
}

func (f SoftwareFinder) IsInstalledAnywhere(configs []Config) (bool, error) {
	for i, config := range configs {
		installed, err := f.IsInstalled(config)
		// Fail silently unless config is the last chance to find the software
		if err != nil && i == len(configs)-1 {
			return false, err
		}

		if installed {
			return true, nil
		}
	}

	return false, nil
}

func (f SoftwareFinder) IsInstalled(config Config) (bool, error) {
	switch config.ForType {
	case CustomPathFinder:
		return f.isInstalledAccordingToCustomPath(config)
	default:
		return f.isInstalledAccordingToRegistry(config)
	}
}

func (f SoftwareFinder) isInstalledAccordingToRegistry(config Config) (bool, error) {
	_, err := f.getInstallDirFromRegistry(config)

	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (f SoftwareFinder) isInstalledAccordingToCustomPath(config Config) (bool, error) {
	return internal.IsValidDirPath(config.CustomInstallPath)
}

func (f SoftwareFinder) GetInstallDirFromSomewhere(configs []Config) (string, error) {
	for i, config := range configs {
		installDir, err := f.GetInstallDir(config)
		// Fail silently unless config is the last chance to find the software
		if err != nil && i == len(configs)-1 {
			return "", err
		}

		if installDir != "" {
			return installDir, nil
		}
	}

	return "", fmt.Errorf("install path could not be determined")
}

func (f SoftwareFinder) GetInstallDir(config Config) (string, error) {
	path, err := f.getInstallDir(config)
	if err != nil {
		return "", err
	}

	// Make sure to return a directory path, even if registry/user-provided config point to a file
	// Also validates paths found in registry also exist on disk
	pathCandidates := []string{path, filepath.Dir(path)}
	for _, candidate := range pathCandidates {
		isDir, err := internal.IsValidDirPath(candidate)
		if err != nil {
			return "", err
		}

		if isDir {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("failed to determine install path based on received path: %s", path)
}

func (f SoftwareFinder) getInstallDir(config Config) (string, error) {
	switch config.ForType {
	case CustomPathFinder:
		return config.CustomInstallPath, nil
	default:
		return f.getInstallDirFromRegistry(config)
	}
}

func (f SoftwareFinder) getInstallDirFromRegistry(config Config) (string, error) {
	return f.repository.GetStringValue(registry.LOCAL_MACHINE, config.RegistryPath, config.RegistryValueName)
}
