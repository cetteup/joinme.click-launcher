package platforms

import (
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var EaClient = domain.PlatformClient{
	Platform: domain.PlatformEaApp,
	FinderConfig: software_finder.Config{
		ForType:           software_finder.RegistryFinder,
		RegistryKey:       software_finder.RegistryKeyLocalMachine,
		RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Desktop",
		RegistryValueName: "ClientPath",
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "EALauncher.exe",
	},
}
