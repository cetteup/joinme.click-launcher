//go:build unit

package registry_repository

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func TestRegistryRepository_GetStringValue(t *testing.T) {
	registryRepository := NewRegistryRepository()

	t.Run("successfully retrieves string value", func(t *testing.T) {
		// GIVEN
		key := registry.LOCAL_MACHINE
		path := "SYSTEM\\CurrentControlSet\\Control\\ComputerName\\ComputerName"
		valueName := "ComputerName"

		expectedValue, err := os.Hostname()
		require.NoError(t, err)

		// WHEN
		value, err := registryRepository.GetStringValue(key, path, valueName)

		// THEN
		require.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("error for non-existing path", func(t *testing.T) {
		// GIVEN
		key := registry.LOCAL_MACHINE
		path := "this-does-not-exist"
		valueName := "ComputerName"

		// WHEN
		_, err := registryRepository.GetStringValue(key, path, valueName)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("error for non-existing value name", func(t *testing.T) {
		// GIVEN
		key := registry.LOCAL_MACHINE
		path := "SYSTEM\\CurrentControlSet\\Control\\ComputerName\\ComputerName"
		valueName := "this-does-not-exist"

		// WHEN
		_, err := registryRepository.GetStringValue(key, path, valueName)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("error for non-string value", func(t *testing.T) {
		// GIVEN
		key := registry.LOCAL_MACHINE
		path := "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion"
		valueName := "InstallTime"

		// WHEN
		_, err := registryRepository.GetStringValue(key, path, valueName)

		// THEN
		require.ErrorContains(t, err, "unexpected key value type")
	})
}

func TestRegistryRepository_DeleteValue(t *testing.T) {
	registryRepository := NewRegistryRepository()
	rand.Seed(time.Now().UnixNano())

	t.Run("successfully deletes value", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "Environment"
		valueName := fmt.Sprintf("some-test-valueName-%d", rand.Int()%512)
		value := fmt.Sprintf("some-test-valueName-%d", rand.Int()%512)

		err := registryRepository.SetStringValue(key, path, valueName, value)
		require.NoError(t, err)

		// WHEN
		err = registryRepository.DeleteValue(key, path, valueName)

		// THEN
		require.NoError(t, err)
		_, err = registryRepository.GetStringValue(key, path, valueName)
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("error for non-existing path", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "this-does-not-exist"
		valueName := fmt.Sprintf("some-test-valueName-%d", rand.Int()%512)

		// WHEN
		err := registryRepository.DeleteValue(key, path, valueName)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("error for non-existing value name", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "Environment"
		valueName := "this-does-not-exist"

		// WHEN
		err := registryRepository.DeleteValue(key, path, valueName)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})
}

func TestRegistryRepository_SetStringValue(t *testing.T) {
	registryRepository := NewRegistryRepository()
	rand.Seed(time.Now().UnixNano())

	t.Run("successfully sets string value", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "Environment"
		valueName := fmt.Sprintf("some-test-valueName-%d", rand.Int()%512)
		value := fmt.Sprintf("some-test-value-%d", rand.Int()%512)

		// WHEN
		err := registryRepository.SetStringValue(key, path, valueName, value)

		t.Cleanup(func() {
			_ = registryRepository.DeleteValue(key, path, valueName)
		})

		// THEN
		require.NoError(t, err)
		actual, err := registryRepository.GetStringValue(key, path, valueName)
		require.NoError(t, err)
		assert.Equal(t, value, actual)
	})

	t.Run("errors for non-existing path", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "this-does-not-exist"
		valueName := fmt.Sprintf("some-test-valueName-%d", rand.Int()%512)
		value := fmt.Sprintf("some-test-value-%d", rand.Int()%512)

		// WHEN
		err := registryRepository.SetStringValue(key, path, valueName, value)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("errors for value name exceeding max length", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "Environment"
		// reference: https://docs.microsoft.com/en-us/windows/win32/sysinfo/registry-element-size-limits
		valueName := strings.Repeat("f", int(math.Pow(2, 14)))
		value := fmt.Sprintf("some-test-value-%d", rand.Int()%512)

		// WHEN
		err := registryRepository.SetStringValue(key, path, valueName, value)

		// THEN
		require.ErrorIs(t, err, windows.ERROR_INVALID_PARAMETER)
	})
}

func TestRegistryRepository_CreateKey(t *testing.T) {
	registryRepository := NewRegistryRepository()
	rand.Seed(time.Now().UnixNano())

	t.Run("successfully creates key", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := fmt.Sprintf("SOFTWARE\\some-test-key-%d", rand.Int()%512)

		// WHEN
		err := registryRepository.CreateKey(key, path)

		t.Cleanup(func() {
			_ = registryRepository.DeleteKey(key, path)
		})

		// THEN
		require.NoError(t, err)
		err = registryRepository.SetStringValue(key, path, "", "")
		require.NoError(t, err)
	})

	t.Run("error for key exceeding max length", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		// reference: https://docs.microsoft.com/en-us/windows/win32/sysinfo/registry-element-size-limits
		path := strings.Repeat("f", int(math.Pow(2, 8))+1)

		// WHEN
		err := registryRepository.CreateKey(key, path)

		// THEN
		require.ErrorIs(t, err, windows.ERROR_INVALID_PARAMETER)
	})
}

func TestRegistryRepository_DeleteKey(t *testing.T) {
	registryRepository := NewRegistryRepository()
	rand.Seed(time.Now().UnixNano())

	t.Run("successfully deletes key", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := fmt.Sprintf("SOFTWARE\\some-test-key-%d", rand.Int()%512)

		err := registryRepository.CreateKey(key, path)

		// WHEN
		err = registryRepository.DeleteKey(key, path)

		// THEN
		require.NoError(t, err)
		err = registryRepository.SetStringValue(key, path, "", "")
		require.ErrorIs(t, err, registry.ErrNotExist)
	})

	t.Run("error for non-existing key", func(t *testing.T) {
		// GIVEN
		key := registry.CURRENT_USER
		path := "SOFTWARE\\this-does-not-exist"

		// WHEN
		err := registryRepository.DeleteKey(key, path)

		// THEN
		require.ErrorIs(t, err, registry.ErrNotExist)
	})
}
