package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Cod2 = title.GameTitle{
	ProtocolScheme: "cod2",
	GameLabel:      "Call of Duty 2",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 2",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "CoD2MP_s.exe",
		CloseBeforeLaunch: true,
	},
	CmdBuilder: PlusConnectCmdBuilder,
}
