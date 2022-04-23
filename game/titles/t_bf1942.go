package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

var Bf1942 = title.GameTitle{
	ProtocolScheme: "bf1942",
	GameLabel:      "Battlefield 1942",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs:    []string{"+restart", "1"},
		ExecutablePath: "BF1942.exe",
	},
	CmdBuilder: bf1942CmdBuilder,
}

var bf1942CmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	args := []string{
		"+joinServer", host,
		"+port", port,
	}

	return args, nil
}

var Bf1942RoadToRome = title.GameTitle{
	ProtocolScheme: "bf1942rtr",
	GameLabel:      "Battlefield 1942: The Road to Rome",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack1",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs: []string{
			"+restart", "1",
			"+game", "XPack1",
		},
		ExecutablePath: "BF1942.exe",
	},
	CmdBuilder: bf1942CmdBuilder,
}

var Bf1942SecretWeaponsOfWW2 = title.GameTitle{
	ProtocolScheme: "bf1942sw",
	GameLabel:      "Battlefield 1942: Secret Weapons of WWII",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack2",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs: []string{
			"+restart", "1",
			"+game", "XPack2",
		},
		ExecutablePath: "BF1942.exe",
	},
	CmdBuilder: bf1942CmdBuilder,
}
