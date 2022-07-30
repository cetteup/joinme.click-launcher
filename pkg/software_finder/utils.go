package software_finder

import (
	"errors"
	"os"
)

type PathType int

const (
	PathTypeFile = iota
	PathTypeDir
)

func PathExistsAndIsType(path string, expect PathType) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return expect == PathTypeFile && !info.IsDir() || expect == PathTypeDir && info.IsDir(), nil
}

func DirExists(path string) (bool, error) {
	existsAndIsType, err := PathExistsAndIsType(path, PathTypeDir)
	if err != nil {
		return false, err
	}
	return existsAndIsType, nil
}
