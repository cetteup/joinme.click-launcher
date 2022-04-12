package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var FearSec2 = title.GameTitle{
	ProtocolScheme: "fearsec2",
	GameLabel:      "F.E.A.R. Combat (SEC2)",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\FEAR-Community.org\\FEAR Combat (SEC2)",
			RegistryValueName: "Path",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "FEARMP.exe",
	},
	CmdBuilder: func(scheme string, ip string, port string) ([]string, error) {
		return []string{"+join", fmt.Sprintf("%s:%s", ip, port)}, nil
	},
}
