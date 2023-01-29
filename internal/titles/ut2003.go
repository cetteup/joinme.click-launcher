package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var UT2003 = domain.GameTitle{
	Name:           "Unreal Tournament 2003",
	ProtocolScheme: "ut2003",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Unreal Technology\\Installed Apps\\UT2003",
			RegistryValueName: "Folder",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "UT2003.exe",
		ExecutablePath: "System",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   internal.MakeSimpleCmdBuilder(),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}
