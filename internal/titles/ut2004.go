package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var UT2004 = domain.GameTitle{
	Name:           "Unreal Tournament 2004",
	ProtocolScheme: "ut2004",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\Unreal Technology\\Installed Apps\\UT2004",
			RegistryValueName: "Folder",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyCurrentUser, // When installed via Steam, key is CurrentUser instead of LocalMachine
			RegistryPath:      "SOFTWARE\\Unreal Technology\\Installed Apps\\UT2004",
			RegistryValueName: "Folder",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Unreal Technology\\Installed Apps\\UT2004", // Disk version uses WOW6432Node
			RegistryValueName: "Folder",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "UT2004.exe",
		ExecutablePath: "System",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator{},
	CmdBuilder:   internal.MakeSimpleCmdBuilder(),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}
