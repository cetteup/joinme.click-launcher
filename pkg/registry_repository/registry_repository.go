package registry_repository

import (
	"golang.org/x/sys/windows/registry"
)

type RegistryRepository struct {
}

func New() *RegistryRepository {
	return &RegistryRepository{}
}

func (r *RegistryRepository) OpenKey(k registry.Key, path string, access uint32, cb func(key registry.Key) error) error {
	key, err := registry.OpenKey(k, path, access)
	if err != nil {
		return err
	}
	defer func(key registry.Key) {
		_ = key.Close()
	}(key)

	return cb(key)
}

func (r *RegistryRepository) GetStringValue(k registry.Key, path string, valueName string) (string, error) {
	var value string
	err := r.OpenKey(k, path, registry.QUERY_VALUE, func(key registry.Key) error {
		var err error
		value, _, err = key.GetStringValue(valueName)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return value, nil
}

func (r *RegistryRepository) SetStringValue(k registry.Key, path string, valueName string, value string) error {
	return r.OpenKey(k, path, registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		return key.SetStringValue(valueName, value)
	})
}

func (r *RegistryRepository) DeleteValue(k registry.Key, path string, valueName string) error {
	return r.OpenKey(k, path, registry.QUERY_VALUE|registry.SET_VALUE, func(key registry.Key) error {
		return key.DeleteValue(valueName)
	})
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
