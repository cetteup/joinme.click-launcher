package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

var Swat4 = title.GameTitle{
	ProtocolScheme: "swat4",
	GameLabel:      "SWAT 4",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\GOG.com\\Games\\1409964317",
			RegistryValueName: "PATH",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "Swat4.exe",
		ExecutablePath:    "Content\\System",
		StartIn:           launcher.BinaryDir,
		CloseBeforeLaunch: true,
	},
	CmdBuilder: swat4CmdBuilder,
}

var swat4CmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	return []string{fmt.Sprintf("%s:%s", host, port)}, nil
}

var Swat4X = title.GameTitle{
	ProtocolScheme: "swat4x",
	GameLabel:      "SWAT 4: The Stetchkov Syndicate",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4 - THE STETCHKOV SYNDICATE",
			RegistryValueName: "InstallPath",
		},
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\GOG.com\\Games\\1409964317",
			RegistryValueName: "PATH",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "Swat4X.exe",
		ExecutablePath:    "ContentExpansion\\System",
		StartIn:           launcher.BinaryDir,
		CloseBeforeLaunch: true,
	},
	CmdBuilder: swat4CmdBuilder,
}
