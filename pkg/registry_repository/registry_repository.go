package registry_repository

import (
	"golang.org/x/sys/windows/registry"
)

type RegistryRepository struct {
}

func New() *RegistryRepository {
	return &RegistryRepository{}
}

func (r *RegistryRepository) GetStringValue(k registry.Key, path string, valueName string) (string, error) {
	key, err := registry.OpenKey(k, path, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	value, _, err := key.GetStringValue(valueName)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RegistryRepository) SetStringValue(k registry.Key, path string, valueName string, value string) error {
	key, err := registry.OpenKey(k, path, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.SetStringValue(valueName, value)
}

func (r *RegistryRepository) DeleteValue(k registry.Key, path string, valueName string) error {
	key, err := registry.OpenKey(k, path, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return key.DeleteValue(valueName)
}

func (r *RegistryRepository) CreateKey(k registry.Key, path string) error {
	key, _, err := registry.CreateKey(k, path, registry.QUERY_VALUE|registry.SET_VALUE)
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)
	return err
}

func (r *RegistryRepository) DeleteKey(k registry.Key, path string) error {
	return registry.DeleteKey(k, path)
}
