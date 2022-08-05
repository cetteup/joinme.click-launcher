package software_finder

import (
	"errors"
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

type FinderType string
type PathType int

const (
	RegistryFinder FinderType = "RegistryFinder"
	PathFinder     FinderType = "PathFinder"

	PathTypeFile = iota
	PathTypeDir
)

type registryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type fileRepository interface {
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
}

type Config struct {
	ForType           FinderType
	RegistryPath      string
	RegistryValueName string
	InstallPath       string
	PathType          PathType
}

type SoftwareFinder struct {
	registryRepository registryRepository
	fileRepository     fileRepository
}

func NewSoftwareFinder(repository registryRepository, fileRepository fileRepository) *SoftwareFinder {
	return &SoftwareFinder{
		registryRepository: repository,
		fileRepository:     fileRepository,
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
	case PathFinder:
		return f.isInstalledAccordingToPath(config)
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

func (f SoftwareFinder) isInstalledAccordingToPath(config Config) (bool, error) {
	switch config.PathType {
	case PathTypeFile:
		return f.fileRepository.FileExists(config.InstallPath)
	case PathTypeDir:
		return f.fileRepository.DirExists(config.InstallPath)
	default:
		return false, fmt.Errorf("unsupported path type: %d", config.PathType)
	}
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
		isDir, err := f.fileRepository.DirExists(candidate)
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
	case PathFinder:
		return f.getInstallDirFromPath(config)
	default:
		return f.getInstallDirFromRegistry(config)
	}
}

func (f SoftwareFinder) getInstallDirFromPath(config Config) (string, error) {
	switch config.PathType {
	case PathTypeFile:
		return filepath.Dir(config.InstallPath), nil
	case PathTypeDir:
		return config.InstallPath, nil
	default:
		return "", fmt.Errorf("unsupported path type: %d", config.PathType)
	}
}

func (f SoftwareFinder) getInstallDirFromRegistry(config Config) (string, error) {
	return f.registryRepository.GetStringValue(registry.LOCAL_MACHINE, config.RegistryPath, config.RegistryValueName)
}
