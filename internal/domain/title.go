package domain

import (
	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

type GameTitle struct {
	ProtocolScheme string
	GameLabel      string
	PlatformClient *Client
	FinderConfigs  []software_finder.Config
	LauncherConfig game_launcher.Config
	URLValidator   game_launcher.URLValidator
	CmdBuilder     game_launcher.CommandBuilder
}

func (t *GameTitle) AddCustomConfig(config internal.CustomLauncherConfig) {
	if config.HasExecutableName() {
		t.LauncherConfig.ExecutableName = config.ExecutableName
	}

	if config.HasExecutablePath() {
		t.LauncherConfig.ExecutablePath = config.ExecutablePath
	}

	if config.HasInstallPath() {
		// Prepend custom path based finder in order search any custom paths first
		t.FinderConfigs = append([]software_finder.Config{
			{
				ForType:           software_finder.CustomPathFinder,
				CustomInstallPath: config.InstallPath,
			},
		}, t.FinderConfigs...)
	}

	if config.HasArgs() {
		t.LauncherConfig.DefaultArgs = append(t.LauncherConfig.DefaultArgs, config.Args...)
	}
}

func (t GameTitle) RequiresPlatformClient() bool {
	return t.PlatformClient != nil
}
