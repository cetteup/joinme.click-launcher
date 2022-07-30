package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bf1942ModPathTemplate       = "Mods\\%s\\lexiconAll.dat"
	bf1942ModRoadToRome         = "Xpack1"
	bf1942ModSecretWeaponsOfWW2 = "Xpack2"
	bf1942Mod1918               = "bf1918"
	bf1942ModDCFinal            = "DC_Final"
	bf1942ModDesertCombat       = "DesertCombat"
	bf1942ModPirates            = "Pirates"
)

var Bf1942 = domain.GameTitle{
	ProtocolScheme: "bf1942",
	GameLabel:      "Battlefield 1942",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs:       []string{"+restart", "1"},
		ExecutableName:    "BF1942.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   bf1942CmdBuilder,
}

var bf1942CmdBuilder game_launcher.CommandBuilder = func(u *url.URL, config game_launcher.Config) ([]string, error) {
	args := []string{
		"+joinServer", u.Hostname(),
		"+port", u.Port(),
	}

	query := u.Query()
	if internal.QueryHasMod(query) {
		mod, err := internal.GetValidModFromQuery(
			query,
			config.InstallPath,
			bf1942ModPathTemplate,
			software_finder.PathTypeFile,
			bf1942ModRoadToRome, bf1942ModSecretWeaponsOfWW2, bf1942Mod1918, bf1942ModDesertCombat, bf1942ModDCFinal, bf1942ModPirates)
		if err != nil {
			return nil, err
		}

		args = append(args, "+game", mod)
	}

	return args, nil
}

var Bf1942RoadToRome = domain.GameTitle{
	ProtocolScheme: "bf1942rtr",
	GameLabel:      "Battlefield 1942: The Road to Rome",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack1",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs: []string{
			"+restart", "1",
			"+game", "XPack1",
		},
		ExecutableName:    "BF1942.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   bf1942CmdBuilder,
}

var Bf1942SecretWeaponsOfWW2 = domain.GameTitle{
	ProtocolScheme: "bf1942sw",
	GameLabel:      "Battlefield 1942: Secret Weapons of WWII",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack2",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs: []string{
			"+restart", "1",
			"+game", "XPack2",
		},
		ExecutableName:    "BF1942.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder:   bf1942CmdBuilder,
}
