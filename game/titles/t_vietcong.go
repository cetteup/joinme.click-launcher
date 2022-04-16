package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
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
		ExecutablePath: "vietcong.exe",
	},
	CmdBuilder: func(scheme string, host string, port string, u *url.URL) ([]string, error) {
		return []string{"-ip", host, "-port", port}, nil
	},
}
