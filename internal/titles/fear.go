package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Fear = domain.GameTitle{
	Name:           "F.E.A.R./F.E.A.R. Combat", // Combat is the free to play release of the multiplayer, both are compatible with patch >=1.07
	ProtocolScheme: "fear",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Monolith Productions\\FEAR\\1.00.0000",
			RegistryValueName: "InstallDir",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Monolith Productions\\FEARCombat\\1.00.0000",
			RegistryValueName: "InstallDir",
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
	URLValidator: internal.IPPortURLValidator{},
	CmdBuilder:   internal.MakeSimpleCmdBuilder("+join"),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}

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
	URLValidator: internal.IPPortURLValidator{},
	CmdBuilder:   internal.MakeSimpleCmdBuilder("+join"),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}
