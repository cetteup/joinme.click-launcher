package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Cod4 = title.GameTitle{
	ProtocolScheme: "cod4",
	GameLabel:      "Call of Duty 4: Modern Warfare",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 4",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "iw3mp.exe",
	},
	CmdBuilder: PlusConnectCmdBuilder,
}
