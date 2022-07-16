package title

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/platform"
	"github.com/cetteup/joinme.click-launcher/internal"
)

type GameTitle struct {
	ProtocolScheme string
	GameLabel      string
	RequiresPort   bool
	PlatformClient *platform.Client
	FinderConfigs  []finder.Config
	LauncherConfig launcher.Config
	CmdBuilder     launcher.CommandBuilder
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
		t.FinderConfigs = append([]finder.Config{
			{
				ForType:           finder.CustomPathFinder,
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
