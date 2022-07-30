package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Swat4 = domain.GameTitle{
	ProtocolScheme: "swat4",
	GameLabel:      "SWAT 4",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\GOG.com\\Games\\1409964317",
			RegistryValueName: "PATH",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "Swat4.exe",
		ExecutablePath:    "Content\\System",
		StartIn:           game_launcher.LaunchDirBinaryDir,
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   swat4CmdBuilder,
}

var swat4CmdBuilder game_launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	return []string{fmt.Sprintf("%s:%s", host, port)}, nil
}

var Swat4X = domain.GameTitle{
	ProtocolScheme: "swat4x",
	GameLabel:      "SWAT 4: The Stetchkov Syndicate",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4 - THE STETCHKOV SYNDICATE",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\GOG.com\\Games\\1409964317",
			RegistryValueName: "PATH",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "Swat4X.exe",
		ExecutablePath:    "ContentExpansion\\System",
		StartIn:           game_launcher.LaunchDirBinaryDir,
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   swat4CmdBuilder,
}
