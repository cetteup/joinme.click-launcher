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
	RequiresPort:   true,
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
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		return []string{"-ip", host, "-port", port}, nil
	},
}
