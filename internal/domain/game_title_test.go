//go:build unit

package domain

import (
	"testing"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/stretchr/testify/assert"
)

func TestGameTitle_AddCustomConfig(t *testing.T) {
	type test struct {
		name        string
		givenConfig internal.CustomLauncherConfig
		expect      func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle)
	}

	tests := []test{
		{
			name: "successfully adds all attributes from custom config",
			givenConfig: internal.CustomLauncherConfig{
				ExecutableName: "custom.exe",
				ExecutablePath: "custom-bin",
				InstallPath:    "C:\\custom",
				Args:           []string{"+custom"},
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				for _, argument := range givenConfig.Args {
					assert.Contains(t, title.LauncherConfig.DefaultArgs, argument)
				}
				assert.Contains(t, title.FinderConfigs, software_finder.Config{
					ForType:     software_finder.PathFinder,
					InstallPath: givenConfig.InstallPath,
					PathType:    software_finder.PathTypeDir,
				})
			},
		},
		{
			name: "successfully adds executable name only",
			givenConfig: internal.CustomLauncherConfig{
				ExecutableName: "custom.exe",
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.LauncherConfig.DefaultArgs, title.LauncherConfig.DefaultArgs)
				assert.Equal(t, givenTitle.FinderConfigs, title.FinderConfigs)
			},
		},
		{
			name: "successfully adds executable path only",
			givenConfig: internal.CustomLauncherConfig{
				ExecutablePath: "custom-bin",
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.DefaultArgs, title.LauncherConfig.DefaultArgs)
				assert.Equal(t, givenTitle.FinderConfigs, title.FinderConfigs)
			},
		},
		{
			name: "successfully adds install path finder config only",
			givenConfig: internal.CustomLauncherConfig{
				InstallPath: "C:\\custom",
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenTitle.LauncherConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.LauncherConfig.DefaultArgs, title.LauncherConfig.DefaultArgs)
				assert.Contains(t, title.FinderConfigs, software_finder.Config{
					ForType:     software_finder.PathFinder,
					InstallPath: givenConfig.InstallPath,
					PathType:    software_finder.PathTypeDir,
				})
			},
		},
		{
			name: "successfully adds arguments only",
			givenConfig: internal.CustomLauncherConfig{
				Args: []string{"+custom"},
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				for _, argument := range givenConfig.Args {
					assert.Contains(t, title.LauncherConfig.DefaultArgs, argument)
				}
				assert.Equal(t, givenTitle.LauncherConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.FinderConfigs, title.FinderConfigs)
			},
		},
		{
			name:        "does not change config if custom config is empty",
			givenConfig: internal.CustomLauncherConfig{},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenTitle.LauncherConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.LauncherConfig.DefaultArgs, title.LauncherConfig.DefaultArgs)
				assert.Equal(t, givenTitle.FinderConfigs, title.FinderConfigs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenTitle := GameTitle{
				FinderConfigs: []software_finder.Config{
					{
						ForType:           software_finder.RegistryFinder,
						RegistryPath:      "default",
						RegistryValueName: "default",
					},
				},
				LauncherConfig: game_launcher.Config{
					ExecutableName: "default.exe",
					ExecutablePath: "default-path",
					DefaultArgs:    []string{"+default"},
				},
			}
			// Copy title so we can compare against original
			title := givenTitle

			// WHEN
			title.AddCustomConfig(tt.givenConfig)

			// THEN
			// InstallPath should still be empty (to be set by finder)
			assert.Equal(t, "", title.LauncherConfig.InstallPath)
			// All original finder configs should still be present
			for _, config := range givenTitle.FinderConfigs {
				assert.Contains(t, title.FinderConfigs, config)
			}
			tt.expect(t, givenTitle, tt.givenConfig, title)
		})
	}
}
