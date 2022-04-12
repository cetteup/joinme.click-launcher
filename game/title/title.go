package title

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/internal"
)

type GameTitle struct {
	ProtocolScheme string
	GameLabel      string
	FinderConfigs  []finder.Config
	LauncherConfig launcher.Config
	CmdBuilder     launcher.CommandBuilder
}

func (t GameTitle) AddCustomConfig(config internal.CustomLauncherConfig) {
	if config.HasInstallPath() {
		// Prepend custom path based finder in order search any custom paths first
		t.FinderConfigs = append([]finder.Config{
			{
				ForType:           finder.CustomPathFinder,
				CustomInstallPath: config.InstallPath,
			},
		}, t.FinderConfigs...)
	}

	if config.HasExecutablePath() {
		t.LauncherConfig.ExecutablePath = config.ExecutablePath
	}

	if config.HasArgs() {
		t.LauncherConfig.DefaultArgs = append(t.LauncherConfig.DefaultArgs, config.Args...)
	}
}
