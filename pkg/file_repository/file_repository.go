package file_repository

import (
	"errors"
	"os"
)

type pathType int

const (
	pathTypeFile = iota
	pathTypeDir
)

type FileRepository struct {
}

func New() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) pathExistsAndIsType(path string, expect pathType) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return expect == pathTypeFile && !info.IsDir() || expect == pathTypeDir && info.IsDir(), nil
}

func (r *FileRepository) FileExists(path string) (bool, error) {
	existsAndIsFile, err := r.pathExistsAndIsType(path, pathTypeFile)
	if err != nil {
		return false, err
	}
	return existsAndIsFile, nil
}

func (r *FileRepository) DirExists(path string) (bool, error) {
	existsAndIsDir, err := r.pathExistsAndIsType(path, pathTypeDir)
	if err != nil {
		return false, err
	}
	return existsAndIsDir, nil
}
