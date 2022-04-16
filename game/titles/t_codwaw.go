package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var CodWaw = title.GameTitle{
	ProtocolScheme: "codwaw",
	GameLabel:      "Call of Duty: World at War",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty WAW",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "CoDWaWmp.exe",
	},
	CmdBuilder: PlusConnectCmdBuilder,
}
