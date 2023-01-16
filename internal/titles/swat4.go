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
	Name:           "SWAT 4",
	ProtocolScheme: "swat4",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
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

var swat4CmdBuilder game_launcher.CommandBuilder = func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		return append(config.DefaultArgs, fmt.Sprintf("%s:%s", u.Hostname(), u.Port())), nil
	}
	return nil, nil
}

var Swat4X = domain.GameTitle{
	Name:           "SWAT 4: The Stetchkov Syndicate",
	ProtocolScheme: "swat4x",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4 - THE STETCHKOV SYNDICATE",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
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
