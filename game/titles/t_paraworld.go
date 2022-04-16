package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

var Paraworld = title.GameTitle{
	ProtocolScheme: "paraworld",
	GameLabel:      "ParaWorld",
	RequiresPort:   true,
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
	CmdBuilder: func(scheme string, host string, port string, u *url.URL) ([]string, error) {
		return []string{"-autoconnect", fmt.Sprintf("%s:%s", host, port)}, nil
	},
}
