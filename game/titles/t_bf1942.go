package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

const (
	bf1942ModBasePath           = "Mods"
	bf1942ModRoadToRome         = "Xpack1"
	bf1942ModSecretWeaponsOfWW2 = "Xpack2"
	bf1942Mod1918               = "bf1918"
	bf1942ModDCFinal            = "DC_Final"
	bf1942ModDesertCombat       = "DesertCombat"
	bf1942ModPirates            = "Pirates"
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

	query := u.Query()
	if query != nil && query.Has(UrlQueryKeyMod) {
		mod, err := getValidMod(
			installPath,
			bf1942ModBasePath,
			query.Get(UrlQueryKeyMod),
			bf1942ModRoadToRome, bf1942ModSecretWeaponsOfWW2, bf1942Mod1918, bf1942ModDesertCombat, bf1942ModDCFinal, bf1942ModPirates,
		)
		if err != nil {
			return nil, err
		}

		args = append(args, "+game", mod)
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
