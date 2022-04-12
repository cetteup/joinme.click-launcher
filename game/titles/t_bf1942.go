package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Bf1942 = title.GameTitle{
	ProtocolScheme: "bf1942",
	GameLabel:      "Battlefield 1942",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs:    []string{"+restart", "1"},
		ExecutablePath: "BF1942.exe",
	},
	CmdBuilder: func(scheme string, ip string, port string) ([]string, error) {
		args := []string{
			"+joinServer", ip,
			"+port", port,
		}

		return args, nil
	},
}
