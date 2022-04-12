package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var BfVietnam = title.GameTitle{
	ProtocolScheme: "bfvietnam",
	GameLabel:      "Battlefield Vietnam",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield Vietnam",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs:    []string{"+restart", "1"},
		ExecutablePath: "BfVietnam.exe",
	},
	CmdBuilder: func(scheme string, ip string, port string) ([]string, error) {
		args := []string{
			"+joinServer", ip,
			"+port", port,
		}

		return args, nil
	},
}
