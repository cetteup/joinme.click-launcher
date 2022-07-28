package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Cod4 = title.GameTitle{
	ProtocolScheme: "cod4",
	GameLabel:      "Call of Duty 4: Modern Warfare",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 4",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "iw3mp.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: ipPortURLValidator,
	CmdBuilder:   plusConnectCmdBuilder,
}
