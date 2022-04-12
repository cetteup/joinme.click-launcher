package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Cod2 = title.GameTitle{
	ProtocolScheme: "cod2",
	GameLabel:      "Call of Duty 2",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 2",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "CoD2MP_s.exe",
	},
	CmdBuilder: PlusConnectCmdBuilder,
}
