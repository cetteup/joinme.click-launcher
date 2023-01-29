package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Cod4 = domain.GameTitle{
	Name:           "Call of Duty 4: Modern Warfare",
	ProtocolScheme: "cod4",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 4",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "iw3mp.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   internal.MakeSimpleCmdBuilder(internal.PlusConnectPrefix),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}
