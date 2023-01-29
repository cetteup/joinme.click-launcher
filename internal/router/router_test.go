//go:build unit

package router

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

func TestGameRouter_AddTitle(t *testing.T) {
	t.Run("successfully adds title", func(t *testing.T) {
		// GIVEN
		router, _, _, _ := getRouterWithDependencies(t)
		protocol := "some-game-protocol"
		title := domain.GameTitle{
			Name:           "some-game",
			ProtocolScheme: protocol,
		}

		// WHEN
		router.AddTitle(title)

		// THEN
		assert.Equal(t, title, router.GameTitles[protocol])
	})

	t.Run("custom config is applied to added title", func(t *testing.T) {
		// GIVEN
		router, _, _, _ := getRouterWithDependencies(t)
		protocol := "some-game-protocol"
		title := domain.GameTitle{
			Name:           "some-game",
			ProtocolScheme: protocol,
		}
		customLauncherConfig := internal.CustomLauncherConfig{
			ExecutableName: "game.exe",
			ExecutablePath: "bin",
			InstallPath:    "C:\\Games\\SomeGame",
			Args:           []string{"+some-arg"},
		}
		internal.Config.Games = map[string]internal.CustomLauncherConfig{
			protocol: customLauncherConfig,
		}

		// WHEN
		router.AddTitle(title)

		// THEN
		assert.Equal(t, title.Name, router.GameTitles[protocol].Name)
		assert.Equal(t, title.ProtocolScheme, router.GameTitles[protocol].ProtocolScheme)
		assert.Equal(t, customLauncherConfig.ExecutableName, router.GameTitles[protocol].LauncherConfig.ExecutableName)
		assert.Equal(t, customLauncherConfig.ExecutablePath, router.GameTitles[protocol].LauncherConfig.ExecutablePath)
		assert.Contains(t, router.GameTitles[protocol].FinderConfigs, software_finder.Config{
			ForType:     software_finder.PathFinder,
			InstallPath: customLauncherConfig.InstallPath,
			PathType:    software_finder.PathTypeDir,
		})
		for _, argument := range customLauncherConfig.Args {
			assert.Contains(t, router.GameTitles[protocol].LauncherConfig.DefaultArgs, argument)
		}
	})
}

func TestGameRouter_RegisterHandlers(t *testing.T) {
	t.Run("successfully registers handlers", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		handlerCommand, err := router.getHandlerCommand()
		require.NoError(t, err)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return("", registry.ErrNotExist)
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(regValueNameDefault), gomock.Eq(fmt.Sprintf("URL:%s protocol", title.Name)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(regValueNameURLProtocol), gomock.Eq(""))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen, regPathCommand})))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen, regPathCommand})), gomock.Eq(regValueNameDefault), gomock.Eq(handlerCommand))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           true,
			PlatformClientInstalled: false,
			PreviouslyRegistered:    false,
			Registered:              true,
			Error:                   nil,
		}, result[0])
	})

	t.Run("successfully updates handler command", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		handlerCommand, err := router.getHandlerCommand()
		require.NoError(t, err)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return("not-a-handler-command", nil)
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(regValueNameDefault), gomock.Eq(fmt.Sprintf("URL:%s protocol", title.Name)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(regValueNameURLProtocol), gomock.Eq(""))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen, regPathCommand})))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{regPathShell, regPathOpen, regPathCommand})), gomock.Eq(regValueNameDefault), gomock.Eq(handlerCommand))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           true,
			PlatformClientInstalled: false,
			PreviouslyRegistered:    false,
			Registered:              true,
			Error:                   nil,
		}, result[0])
	})

	t.Run("checks if required platform client is installed", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		handlerCommand, err := router.getHandlerCommand()
		require.NoError(t, err)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockFinder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(handlerCommand, nil)

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           true,
			PlatformClientInstalled: true,
			PreviouslyRegistered:    true,
			Registered:              false,
			Error:                   nil,
		}, result[0])
	})

	t.Run("skips game if handler is already registered", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		handlerCommand, err := router.getHandlerCommand()
		require.NoError(t, err)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return(handlerCommand, nil)

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           true,
			PlatformClientInstalled: false,
			PreviouslyRegistered:    true,
			Registered:              false,
			Error:                   nil,
		}, result[0])
	})

	t.Run("skips game if not installed", func(t *testing.T) {
		// GIVEN
		router, _, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(false, nil)

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           false,
			PlatformClientInstalled: false,
			PreviouslyRegistered:    false,
			Registered:              false,
			Error:                   nil,
		}, result[0])
	})

	t.Run("skips game if required platform client is not installed", func(t *testing.T) {
		// GIVEN
		router, _, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockFinder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(false, nil)

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, handlerRegistrationResult{
			Title:                   title,
			GameInstalled:           true,
			PlatformClientInstalled: false,
			PreviouslyRegistered:    false,
			Registered:              false,
			Error:                   nil,
		}, result[0])
	})

	t.Run("error if finder encounters an error checking for game", func(t *testing.T) {
		// GIVEN
		router, _, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(false, fmt.Errorf("some-error"))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, title, result[0].Title)
		require.ErrorContains(t, result[0].Error, "failed to determine whether game is installed")
	})

	t.Run("error if finder encounters an error checking for platform client", func(t *testing.T) {
		// GIVEN
		router, _, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockFinder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(false, fmt.Errorf("some-error"))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, title, result[0].Title)
		require.ErrorContains(t, result[0].Error, fmt.Sprintf("failed to determine whether required platform (%s) is installed", title.PlatformClient.Platform))
	})

	t.Run("error if handler registration check fails", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return("", fmt.Errorf("some-error"))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, title, result[0].Title)
		require.ErrorContains(t, result[0].Error, "failed to determine whether handler is registered")
	})

	t.Run("error if handler registration fails", func(t *testing.T) {
		// GIVEN
		router, mockRepository, mockFinder, _ := getRouterWithDependencies(t)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\some-game",
					RegistryValueName: "some-value-name",
				},
			},
		}
		router.AddTitle(title)

		mockFinder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
		mockRepository.EXPECT().GetStringValue(gomock.Any(), gomock.Any(), gomock.Any()).Return("", registry.ErrNotExist)
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil))).Return(fmt.Errorf("some-error"))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, title, result[0].Title)
		require.ErrorContains(t, result[0].Error, "failed to register as URL protocol handler")
	})
}

func TestGameRouter_RunURL(t *testing.T) {
	type test struct {
		name                string
		givenTitle          *domain.GameTitle
		givenCommandLineURL string
		expect              func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher)
		wantTitle           *domain.GameTitle
		wantErrContains     string
	}

	tests := []test{
		{
			name:                "successfully launches game and joins server",
			givenCommandLineURL: "bf2://127.0.0.1:16567",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				gameInstallPath := "C:\\Games\\BF2"
				finder.EXPECT().GetInstallDirFromSomewhere(gomock.Eq(title.FinderConfigs)).Return(gameInstallPath, nil)
				finalLaunchConfig := title.LauncherConfig
				finalLaunchConfig.InstallPath = gameInstallPath
				launcher.EXPECT().StartGame(
					gomock.Eq(&url.URL{
						Scheme: "bf2",
						Host:   "127.0.0.1:16567",
					}),
					gomock.Eq(finalLaunchConfig),
					gomock.Eq(game_launcher.LaunchTypeLaunchAndJoin),
					gomock.Any(),
					gomock.Any(),
				)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "",
		},
		{
			name:                "successfully launches game with platform client and joins server",
			givenCommandLineURL: "bf4://1234567890",
			givenTitle:          &titles.Bf4,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				platformClientInstallPath := "C:\\Games\\Origin"
				finder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(true, nil)
				finder.EXPECT().GetInstallDir(gomock.Eq(title.PlatformClient.FinderConfig)).Return(platformClientInstallPath, nil)
				finalLaunchConfig := title.LauncherConfig
				finalLaunchConfig.ExecutableName = title.PlatformClient.LauncherConfig.ExecutableName
				finalLaunchConfig.ExecutablePath = title.PlatformClient.LauncherConfig.ExecutablePath
				finalLaunchConfig.InstallPath = platformClientInstallPath
				launcher.EXPECT().StartGame(
					gomock.Eq(&url.URL{
						Scheme: "bf4",
						Host:   "1234567890",
					}),
					gomock.Eq(finalLaunchConfig),
					gomock.Eq(game_launcher.LaunchTypeLaunchAndJoin),
					gomock.Any(),
					gomock.Any(),
				)
			},
			wantTitle:       &titles.Bf4,
			wantErrContains: "",
		},
		{
			name:                "successfully launches game with mod and joins server",
			givenCommandLineURL: "bf1942://127.0.0.1:14567?mod=xpack1",
			givenTitle:          &titles.Bf1942,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				gameInstallPath := "C:\\Games\\BF1942"
				finder.EXPECT().GetInstallDirFromSomewhere(gomock.Eq(title.FinderConfigs)).Return(gameInstallPath, nil)
				modFinderConfig := title.Mods[0].ComputeFinderConfigs(gameInstallPath)
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(modFinderConfig)).Return(true, nil)
				finder.EXPECT().GetInstallDirFromSomewhere(gomock.Eq(title.FinderConfigs)).Return(gameInstallPath, nil)
				finalLaunchConfig := title.LauncherConfig
				finalLaunchConfig.InstallPath = gameInstallPath
				launcher.EXPECT().StartGame(
					gomock.Eq(&url.URL{
						Scheme:   "bf1942",
						Host:     "127.0.0.1:14567",
						RawQuery: "mod=xpack1",
					}),
					gomock.Eq(finalLaunchConfig),
					gomock.Eq(game_launcher.LaunchTypeLaunchAndJoin),
					gomock.Any(),
					gomock.Any(),
				)
			},
			wantTitle:       &titles.Bf1942,
			wantErrContains: "",
		},
		{
			name:                "successfully launches game via action URL",
			givenCommandLineURL: "bf2://act/launch",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				gameInstallPath := "C:\\Games\\BF2"
				finder.EXPECT().GetInstallDirFromSomewhere(gomock.Eq(title.FinderConfigs)).Return(gameInstallPath, nil)
				finalLaunchConfig := title.LauncherConfig
				finalLaunchConfig.InstallPath = gameInstallPath
				launcher.EXPECT().StartGame(
					gomock.Eq(&url.URL{
						Scheme: "bf2",
						Host:   "act",
						Path:   "/launch",
					}),
					gomock.Eq(finalLaunchConfig),
					gomock.Eq(game_launcher.LaunchTypeLaunchOnly),
					gomock.Any(),
					gomock.Any(),
				)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "",
		},
		{
			name:                "error for unsupported game",
			givenCommandLineURL: "not-a-supported-game://127.0.0.1:16567",
			expect:              func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {},
			wantErrContains:     "game not supported",
		},
		{
			name:                "error for unsupported mod",
			givenCommandLineURL: "bf2://127.0.0.1:16567?mod=not-a-supported-mod",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "mod not supported",
		},
		{
			name:                "error for unsupported action",
			givenCommandLineURL: "bf2://act/not-a-supported-action",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "action not supported",
		},
		{
			name:                "error for non-installed game",
			givenCommandLineURL: "bf2://127.0.0.1:16567",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(false, nil)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "game not installed",
		},
		{
			name:                "error for non-installed platform client",
			givenCommandLineURL: "bf4://1234567890",
			givenTitle:          &titles.Bf4,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				finder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(false, nil)
			},
			wantTitle:       &titles.Bf4,
			wantErrContains: "required platform client not installed",
		},
		{
			name:                "error for non-installed mod",
			givenCommandLineURL: "bf2://127.0.0.1:16567?mod=xpack",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				gameInstallPath := "C:\\Games\\BF2"
				finder.EXPECT().GetInstallDirFromSomewhere(gomock.Eq(title.FinderConfigs)).Return(gameInstallPath, nil)
				modFinderConfigs := title.Mods[0].ComputeFinderConfigs(gameInstallPath)
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(modFinderConfigs)).Return(false, nil)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "mod not installed",
		},
		{
			name:                "error for non-parseable URL",
			givenCommandLineURL: "://",
			expect:              func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {},
			wantErrContains:     "missing protocol scheme",
		},
		{
			name:                "error for invalid ip:port URL",
			givenCommandLineURL: "bf2://127.0.0.1",
			givenTitle:          &titles.Bf2,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
			},
			wantTitle:       &titles.Bf2,
			wantErrContains: "port is missing from url",
		},
		{
			name:                "error for invalid gameid URL",
			givenCommandLineURL: "bf4://not-a-game-id",
			givenTitle:          &titles.Bf4,
			expect: func(title *domain.GameTitle, finder *MockGameFinder, launcher *MockGameLauncher) {
				finder.EXPECT().IsInstalledAnywhere(gomock.Eq(title.FinderConfigs)).Return(true, nil)
				finder.EXPECT().IsInstalled(gomock.Eq(title.PlatformClient.FinderConfig)).Return(true, nil)
			},
			wantTitle:       &titles.Bf4,
			wantErrContains: "url hostname is not a valid game id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GIVEN
			router, _, mockFinder, mockLauncher := getRouterWithDependencies(t)
			if tt.givenTitle != nil {
				router.AddTitle(*tt.givenTitle)
			}

			// EXPECT
			tt.expect(tt.givenTitle, mockFinder, mockLauncher)

			// WHEN
			title, err := router.RunURL(tt.givenCommandLineURL)
			if tt.wantTitle != nil {
				assert.Equal(t, tt.wantTitle.ProtocolScheme, title.ProtocolScheme)
			} else {
				assert.Nil(t, title)
			}
			if tt.wantErrContains != "" {
				require.ErrorContains(t, err, tt.wantErrContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func getRouterWithDependencies(t *testing.T) (*GameRouter, *MockRegistryRepository, *MockGameFinder, *MockGameLauncher) {
	ctrl := gomock.NewController(t)
	mockRepository := NewMockRegistryRepository(ctrl)
	mockFinder := NewMockGameFinder(ctrl)
	mockLauncher := NewMockGameLauncher(ctrl)
	return New(mockRepository, mockFinder, mockLauncher), mockRepository, mockFinder, mockLauncher
}
