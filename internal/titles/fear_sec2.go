package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var FearSec2 = domain.GameTitle{
	Name:           "F.E.A.R. Combat (SEC2)",
	ProtocolScheme: "fearsec2",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\FEAR-Community.org\\FEAR Combat (SEC2)",
			RegistryValueName: "Path",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "FEARMP.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   internal.MakeSimpleCmdBuilder("+join"),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}
