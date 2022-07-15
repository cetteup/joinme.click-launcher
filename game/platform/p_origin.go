package platform

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var OriginClient = Client{
	Platform: Origin,
	FinderConfig: finder.Config{
		ForType:           finder.RegistryFinder,
		RegistryPath:      "SOFTWARE\\WOW6432Node\\Origin",
		RegistryValueName: "ClientPath",
	},
	LauncherConfig: launcher.Config{
		ExecutableName: "Origin.exe",
	},
}
