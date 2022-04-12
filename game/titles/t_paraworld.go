package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Paraworld = title.GameTitle{
	ProtocolScheme: "paraworld",
	GameLabel:      "ParaWorld",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "bin\\Paraworld.exe",
	},
	CmdBuilder: func(scheme string, ip string, port string) ([]string, error) {
		return []string{"-autoconnect", fmt.Sprintf("%s:%s", ip, port)}, nil
	},
}
