//go:build unit

package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("successfully loads config", func(t *testing.T) {
		// GIVEN
		Config = &config{}
		givenConfig := config{
			DebugLogging: true,
			QuietLaunch:  true,
			Games: map[string]CustomLauncherConfig{
				"some-game": {
					ExecutableName: "game.exe",
					ExecutablePath: "bin",
					InstallPath:    "C:\\Games\\SomeGame",
					Args:           []string{"+fullscreen", "0"},
				},
			},
		}
		content, err := yaml.Marshal(givenConfig)
		require.NoError(t, err)
		configFilePath, err := buildConfigFilePath()
		require.NoError(t, err)
		err = writeConfigFile(configFilePath, content)
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.Remove(configFilePath)
		})

		// WHEN
		err = LoadConfig()
		require.NoError(t, err)

		// THEN
		assert.Equal(t, &givenConfig, Config)
	})

	t.Run("does not error if config file does not exist", func(t *testing.T) {
		// GIVEN
		Config = &config{}

		// WHEN
		err := LoadConfig()

		// THEN
		require.NoError(t, err)
		assert.Equal(t, &config{}, Config)
	})

	t.Run("error if config file contains invalid yaml", func(t *testing.T) {
		// GIVEN
		Config = &config{}
		configFilePath, err := buildConfigFilePath()
		require.NoError(t, err)
		err = writeConfigFile(configFilePath, []byte("this-is-not-valid-yaml"))
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.Remove(configFilePath)
		})

		// WHEN
		err = LoadConfig()

		// THEN
		require.ErrorContains(t, err, "cannot unmarshal")
		assert.Equal(t, &config{}, Config)
	})
}

func TestConfig_GetCustomLauncherConfig(t *testing.T) {
	type test struct {
		name                         string
		givenConfig                  config
		givenGame                    string
		expectedCustomLauncherConfig *CustomLauncherConfig
	}

	tests := []test{
		{
			name: "successfully retrieves custom launcher config",
			givenConfig: config{
				Games: map[string]CustomLauncherConfig{
					"some-game": {
						ExecutableName: "game.exe",
						ExecutablePath: "bin",
						InstallPath:    "C:\\Games\\some-game",
						Args:           []string{"+prio", "high"},
					},
				},
			},
			givenGame: "some-game",
			expectedCustomLauncherConfig: &CustomLauncherConfig{
				ExecutableName: "game.exe",
				ExecutablePath: "bin",
				InstallPath:    "C:\\Games\\some-game",
				Args:           []string{"+prio", "high"},
			},
		},
		{
			name: "nil if config does not contain custom launcher config for game",
			givenConfig: config{
				Games: map[string]CustomLauncherConfig{
					"some-other-game": {
						ExecutableName: "other-game.exe",
						ExecutablePath: "other-bin",
						InstallPath:    "C:\\Games\\some-other-game",
						Args:           []string{"+level", "mp_caretan"},
					},
				},
			},
			givenGame:                    "some-game",
			expectedCustomLauncherConfig: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			customLauncherConfig := tt.givenConfig.GetCustomLauncherConfig(tt.givenGame)

			// THEN
			assert.Equal(t, tt.expectedCustomLauncherConfig, customLauncherConfig)
		})
	}
}

func TestCustomLauncherConfig_HasValues(t *testing.T) {
	type test struct {
		name          string
		givenConfig   *CustomLauncherConfig
		wantHasValues bool
	}

	tests := []test{
		{
			name: "true for config with all attributes",
			givenConfig: &CustomLauncherConfig{
				ExecutableName: "game.exe",
				ExecutablePath: "bin",
				InstallPath:    "C:\\Games\\SomeGame",
				Args:           []string{"+some-arg"},
			},
			wantHasValues: true,
		},
		{
			name: "true for config with executable name only",
			givenConfig: &CustomLauncherConfig{
				ExecutableName: "game.exe",
			},
			wantHasValues: true,
		},
		{
			name: "true for config with executable path only",
			givenConfig: &CustomLauncherConfig{
				ExecutablePath: "bin",
			},
			wantHasValues: true,
		},
		{
			name: "true for config with install path only",
			givenConfig: &CustomLauncherConfig{
				InstallPath: "C:\\Games\\SomeGame",
			},
			wantHasValues: true,
		},
		{
			name: "true for config with args only",
			givenConfig: &CustomLauncherConfig{
				Args: []string{"+some-arg"},
			},
			wantHasValues: true,
		},
		{
			name:          "false for empty config",
			givenConfig:   &CustomLauncherConfig{},
			wantHasValues: false,
		},
		{
			name:          "false for nil config",
			givenConfig:   nil,
			wantHasValues: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			hasValues := tt.givenConfig.HasValues()

			// THEN
			assert.Equal(t, tt.wantHasValues, hasValues)
		})
	}
}

func buildConfigFilePath() (string, error) {
	wd, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(wd), ConfigFilename), nil
}

func writeConfigFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0666)
}
