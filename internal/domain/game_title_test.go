//go:build unit

package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
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
				Hooks: []internal.CustomHookConfig{
					{
						Handler:     "some-handler",
						When:        game_launcher.HookWhenAlways,
						ExitOnError: true,
						Args: map[string]string{
							"some-key": "some-value",
						},
					},
				},
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				assert.Equal(t, givenConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				for _, argument := range givenConfig.Args {
					assert.Contains(t, title.LauncherConfig.DefaultArgs, argument)
				}
				for _, chc := range givenConfig.Hooks {
					assert.Contains(t, title.LauncherConfig.HookConfigs, game_launcher.HookConfig{
						Handler:     chc.Handler,
						When:        chc.When,
						ExitOnError: chc.ExitOnError,
						Args:        chc.Args,
					})
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
				assert.Equal(t, givenTitle.LauncherConfig.HookConfigs, title.LauncherConfig.HookConfigs)
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
				assert.Equal(t, givenTitle.LauncherConfig.HookConfigs, title.LauncherConfig.HookConfigs)
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
				assert.Equal(t, givenTitle.LauncherConfig.HookConfigs, title.LauncherConfig.HookConfigs)
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
				assert.Equal(t, givenTitle.LauncherConfig.HookConfigs, title.LauncherConfig.HookConfigs)
				assert.Equal(t, givenTitle.FinderConfigs, title.FinderConfigs)
			},
		},
		{
			name: "successfully adds hooks only",
			givenConfig: internal.CustomLauncherConfig{
				Hooks: []internal.CustomHookConfig{
					{
						Handler:     "some-handler",
						When:        game_launcher.HookWhenAlways,
						ExitOnError: true,
						Args: map[string]string{
							"some-key": "some-value",
						},
					},
				},
			},
			expect: func(t *testing.T, givenTitle GameTitle, givenConfig internal.CustomLauncherConfig, title GameTitle) {
				for _, chc := range givenConfig.Hooks {
					assert.Contains(t, title.LauncherConfig.HookConfigs, game_launcher.HookConfig{
						Handler:     chc.Handler,
						When:        chc.When,
						ExitOnError: chc.ExitOnError,
						Args:        chc.Args,
					})
				}
				assert.Equal(t, givenTitle.LauncherConfig.ExecutableName, title.LauncherConfig.ExecutableName)
				assert.Equal(t, givenTitle.LauncherConfig.ExecutablePath, title.LauncherConfig.ExecutablePath)
				assert.Equal(t, givenTitle.LauncherConfig.DefaultArgs, title.LauncherConfig.DefaultArgs)
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
						RegistryKey:       software_finder.RegistryKeyLocalMachine,
						RegistryPath:      "default",
						RegistryValueName: "default",
					},
				},
				LauncherConfig: game_launcher.Config{
					ExecutableName: "default.exe",
					ExecutablePath: "default-path",
					DefaultArgs:    []string{"+default"},
					HookConfigs: []game_launcher.HookConfig{
						{
							Handler: "some-default-handler",
							When:    game_launcher.HookWhenAlways,
						},
					},
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

func TestGameTitle_RequiresPlatformClient(t *testing.T) {
	type test struct {
		name                       string
		givenTitle                 GameTitle
		wantRequiresPlatformClient bool
	}

	tests := []test{
		{
			name: "true for game title with platform client",
			givenTitle: GameTitle{
				PlatformClient: &PlatformClient{
					Platform: "some-platform-client",
				},
			},
			wantRequiresPlatformClient: true,
		},
		{
			name:                       "false for game title without platform client",
			givenTitle:                 GameTitle{},
			wantRequiresPlatformClient: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			requiresPlatformClient := tt.givenTitle.RequiresPlatformClient()

			// THEN
			assert.Equal(t, tt.wantRequiresPlatformClient, requiresPlatformClient)
		})
	}
}

func TestGameTitle_GetMod(t *testing.T) {
	type test struct {
		name        string
		givenTitle  GameTitle
		givenSlug   string
		expectedMod *GameMod
	}

	tests := []test{
		{
			name: "successfully returns mod",
			givenTitle: GameTitle{
				Mods: []GameMod{
					{
						Name:          "some-mod",
						Slug:          "some-mod-slug",
						finderConfigs: nil,
					},
				},
			},
			givenSlug: "some-mod-slug",
			expectedMod: &GameMod{
				Name:          "some-mod",
				Slug:          "some-mod-slug",
				finderConfigs: nil,
			},
		},
		{
			name: "ignores case when comparing slugs",
			givenTitle: GameTitle{
				Mods: []GameMod{
					{
						Name:          "some-mod",
						Slug:          "soMe-moD-slUG",
						finderConfigs: nil,
					},
				},
			},
			givenSlug: "SOmE-mOd-SLug",
			expectedMod: &GameMod{
				Name:          "some-mod",
				Slug:          "soMe-moD-slUG",
				finderConfigs: nil,
			},
		},
		{
			name: "nil for unsupported mod",
			givenTitle: GameTitle{
				Mods: []GameMod{
					MakeMod("some-mod", "some-mod-slug", nil),
				},
			},
			givenSlug:   "some-other-mod-slug",
			expectedMod: nil,
		},
		{
			name:        "nil for title without mods",
			givenTitle:  GameTitle{},
			givenSlug:   "some-mod-slug",
			expectedMod: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			mod := tt.givenTitle.GetMod(tt.givenSlug)

			// THEN
			assert.Equal(t, tt.expectedMod, mod)
		})
	}
}

func TestGameTitle_String(t *testing.T) {
	type test struct {
		name           string
		givenTitle     *GameTitle
		expectedString string
	}

	tests := []test{
		{
			name: "returns game name and protocol scheme for given title",
			givenTitle: &GameTitle{
				Name:           "some-game",
				ProtocolScheme: "some-game-protocol",
			},
			expectedString: "some-game (some-game-protocol)",
		},
		{
			name:           "returns nil as string for nil title",
			givenTitle:     nil,
			expectedString: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asString := tt.givenTitle.String()

			// THEN
			assert.Equal(t, tt.expectedString, asString)
		})
	}
}

func TestGameMod_ComputeFinderConfigs(t *testing.T) {
	type test struct {
		name                  string
		givenMod              GameMod
		givenGameInstallPath  string
		expectedFinderConfigs []software_finder.Config
	}

	tests := []test{
		{
			name: "successfully computes finder configs",
			givenMod: GameMod{
				Name: "some-mod",
				Slug: "some-mod-slug",
				finderConfigs: []software_finder.Config{
					{
						ForType:     software_finder.PathFinder,
						InstallPath: "mods\\some-mod\\objects.zip",
						PathType:    software_finder.PathTypeFile,
					},
				},
			},
			givenGameInstallPath: "C:\\Games\\Battlefield",
			expectedFinderConfigs: []software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: "C:\\Games\\Battlefield\\mods\\some-mod\\objects.zip",
					PathType:    software_finder.PathTypeFile,
				},
			},
		},
		{
			name: "ignores finder configs not using path finder",
			givenMod: GameMod{
				Name: "some-mod",
				Slug: "some-mod-slug",
				finderConfigs: []software_finder.Config{
					{
						ForType:           software_finder.RegistryFinder,
						RegistryKey:       software_finder.RegistryKeyLocalMachine,
						RegistryPath:      "some-game\\some-mod",
						RegistryValueName: "InstallPath",
					},
				},
			},
			givenGameInstallPath: "C:\\Games\\Battlefield",
			expectedFinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "some-game\\some-mod",
					RegistryValueName: "InstallPath",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			finderConfigs := tt.givenMod.ComputeFinderConfigs(tt.givenGameInstallPath)

			// THEN
			assert.Equal(t, tt.expectedFinderConfigs, finderConfigs)
		})
	}
}

func TestGameMod_String(t *testing.T) {
	type test struct {
		name           string
		givenMod       *GameMod
		expectedString string
	}

	tests := []test{
		{
			name: "returns mod name and slug for given mod",
			givenMod: &GameMod{
				Name: "some-mod",
				Slug: "some-mod-slug",
			},
			expectedString: "some-mod (some-mod-slug)",
		},
		{
			name:           "returns nil as string for nil title",
			givenMod:       nil,
			expectedString: "nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// WHEN
			asString := tt.givenMod.String()

			// THEN
			assert.Equal(t, tt.expectedString, asString)
		})
	}
}
