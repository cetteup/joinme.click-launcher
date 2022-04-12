package titles

import (
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
		ExecutablePath: "vietcong.exe",
	},
	CmdBuilder: func(scheme string, ip string, port string) ([]string, error) {
		return []string{"-ip", ip, "-port", port}, nil
	},
}
