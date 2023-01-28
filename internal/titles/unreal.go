package titles

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Unreal = domain.GameTitle{
	Name:           "Unreal",
	ProtocolScheme: "unreal",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\Unreal Technology\\Installed Apps\\Unreal Gold",
			RegistryValueName: "folder",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyCurrentUser, // When installed via Steam, key is CurrentUser instead of LocalMachine
			RegistryPath:      "SOFTWARE\\Unreal Technology\\Installed Apps\\Unreal Gold",
			RegistryValueName: "folder",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\App Paths\\Unreal.exe", // TODO Points to the "System" dir for Unreal Gold
			RegistryValueName: "Path",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "Unreal.exe",
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
	CmdBuilder:   internal.PlainCmdBuilder,
	HookHandlers: map[string]game_launcher.HookHandler{
		internal.HookKillProcess: internal.KillProcessHookHandler(true),
	},
}
