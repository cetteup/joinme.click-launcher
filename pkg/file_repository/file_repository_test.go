//go:build unit

package file_repository

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileRepository_FileExists(t *testing.T) {
	fileRepository := New()

	t.Run("true for existing file", func(t *testing.T) {
		// GIVEN
		tempFile, err := createTempFile(t)
		require.NoError(t, err)

		// WHEN
		fileExists, err := fileRepository.FileExists(tempFile)

		// THEN
		require.NoError(t, err)
		assert.True(t, fileExists)
	})

	t.Run("false for mismatched type", func(t *testing.T) {
		// GIVEN
		tempDir, err := createTempDir(t)
		require.NoError(t, err)

		// WHEN
		fileExists, err := fileRepository.FileExists(tempDir)

		// THEN
		require.NoError(t, err)
		assert.False(t, fileExists)

	})

	t.Run("false for non existing file", func(t *testing.T) {
		// GIVEN
		tempFile, err := createTempFile(t)
		require.NoError(t, err)
		err = os.Remove(tempFile)
		require.NoError(t, err)

		// WHEN
		fileExists, err := fileRepository.FileExists(tempFile)

		// THEN
		require.NoError(t, err)
		assert.False(t, fileExists)
	})

	t.Run("error for invalid path", func(t *testing.T) {
		// GIVEN
		path := strings.Repeat("/", 4)

		// WHEN
		fileExists, err := fileRepository.FileExists(path)

		// THEN
		require.Error(t, err)
		assert.False(t, fileExists)
	})
}

func TestFileRepository_DirExists(t *testing.T) {
	fileRepository := New()

	t.Run("true for existing directory", func(t *testing.T) {
		// GIVEN
		tempDir, err := createTempDir(t)
		require.NoError(t, err)

		// WHEN
		dirExists, err := fileRepository.DirExists(filepath.Dir(tempDir))

		// THEN
		require.NoError(t, err)
		assert.True(t, dirExists)
	})

	t.Run("false for mismatched type", func(t *testing.T) {
		tempFile, err := createTempFile(t)
		require.NoError(t, err)

		// WHEN
		dirExists, err := fileRepository.DirExists(tempFile)

		// THEN
		require.NoError(t, err)
		assert.False(t, dirExists)
	})

	t.Run("false for non existing directory", func(t *testing.T) {
		// GIVEN
		tempDir, err := createTempDir(t)
		require.NoError(t, err)
		err = os.Remove(tempDir)
		require.NoError(t, err)

		// WHEN
		dirExists, err := fileRepository.DirExists(tempDir)

		// THEN
		require.NoError(t, err)
		assert.False(t, dirExists)
	})

	t.Run("error for invalid path", func(t *testing.T) {
		// GIVEN
		path := strings.Repeat("/", 4)

		// WHEN
		fileExists, err := fileRepository.DirExists(path)

		// THEN
		require.Error(t, err)
		assert.False(t, fileExists)
	})
}

func createTempFile(t *testing.T) (string, error) {
	f, err := os.CreateTemp("", "test-file")
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})
	return f.Name(), nil
}

func createTempDir(t *testing.T) (string, error) {
	dirPath, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		return "", err
	}
	t.Cleanup(func() {
		_ = os.Remove(dirPath)
	})
	return dirPath, nil
}
