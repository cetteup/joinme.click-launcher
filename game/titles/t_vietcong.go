package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Vietcong = title.GameTitle{
	ProtocolScheme: "vietcong",
	GameLabel:      "Vietcong",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Pterodon\\Vietcong",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "vietcong.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: ipPortURLValidator,
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		return []string{"-ip", host, "-port", port}, nil
	},
}
