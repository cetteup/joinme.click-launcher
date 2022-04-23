package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
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
		ExecutablePath: "Content\\System\\Swat4.exe",
		StartIn:        launcher.BinaryDir,
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
		ExecutablePath: "ContentExpansion\\System\\Swat4X.exe",
		StartIn:        launcher.BinaryDir,
	},
	CmdBuilder: swat4CmdBuilder,
}
