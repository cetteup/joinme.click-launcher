//go:build unit

package router

import (
	"fmt"
	"testing"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/registry"
)

func TestGameRouter_AddTitle(t *testing.T) {
	t.Run("successfully adds title", func(t *testing.T) {
		// GIVEN
		router := getRouterWithDependencies(t)
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
		router := getRouterWithDependencies(t)
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(RegKeyDefault), gomock.Eq(fmt.Sprintf("URL:%s protocol", title.Name)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(RegKeyURLProtocol), gomock.Eq(""))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen, RegPathCommand})))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen, RegPathCommand})), gomock.Eq(RegKeyDefault), gomock.Eq(handlerCommand))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(RegKeyDefault), gomock.Eq(fmt.Sprintf("URL:%s protocol", title.Name)))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, nil)), gomock.Eq(RegKeyURLProtocol), gomock.Eq(""))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen})))
		mockRepository.EXPECT().CreateKey(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen, RegPathCommand})))
		mockRepository.EXPECT().SetStringValue(gomock.Eq(registry.CURRENT_USER), gomock.Eq(router.getUrlHandlerRegistryPath(title, []string{RegPathShell, RegPathOpen, RegPathCommand})), gomock.Eq(RegKeyDefault), gomock.Eq(handlerCommand))

		// WHEN
		result := router.RegisterHandlers()

		// THEN
		assert.Len(t, result, 1)
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		assert.Equal(t, HandlerRegistrationResult{
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			PlatformClient: &domain.PlatformClient{
				Platform: "some-platform",
				FinderConfig: software_finder.Config{
					ForType:           software_finder.RegistryFinder,
					RegistryPath:      "SOFTWARE\\some-client",
					RegistryValueName: "some-value-name",
				},
			},
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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
		ctrl := gomock.NewController(t)
		mockRepository := NewMockregistryRepository(ctrl)
		mockFinder := NewMockgameFinder(ctrl)
		mockLauncher := NewMockgameLauncher(ctrl)
		router := NewGameRouter(mockRepository, mockFinder, mockLauncher)

		title := domain.GameTitle{
			Name:           "some-name",
			ProtocolScheme: "some-protocol",
			FinderConfigs: []software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
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

func getRouterWithDependencies(t *testing.T) *GameRouter {
	ctrl := gomock.NewController(t)
	mockRepository := NewMockregistryRepository(ctrl)
	mockFinder := NewMockgameFinder(ctrl)
	mockLauncher := NewMockgameLauncher(ctrl)
	return NewGameRouter(mockRepository, mockFinder, mockLauncher)
}
