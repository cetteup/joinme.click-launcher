package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Cod = title.GameTitle{
	ProtocolScheme: "cod",
	GameLabel:      "Call of Duty",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "CoDMP.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: ipPortURLValidator,
	CmdBuilder:   plusConnectCmdBuilder,
}

var CodUO = title.GameTitle{
	ProtocolScheme: "coduo",
	GameLabel:      "Call of Duty: United Offensive",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty United Offensive",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "CoDUOMP.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: frostbite3GameIdURLValidator,
	CmdBuilder:   plusConnectCmdBuilder,
}
