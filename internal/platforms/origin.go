package platforms

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var OriginClient = domain.PlatformClient{
	Platform: domain.PlatformOrigin,
	FinderConfig: software_finder.Config{
		ForType:           software_finder.RegistryFinder,
		RegistryPath:      "SOFTWARE\\WOW6432Node\\Origin",
		RegistryValueName: "ClientPath",
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "Origin.exe",
	},
}
