package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

const (
	bf2ProfileFolder    = "Battlefield 2"
	bf2ModBasePath      = "mods"
	bf2ModSpecialForces = "xpack"
	bf2ModAIX2          = "AIX2"
	bf2ModPirates       = "bfp2"
	bf2ModPoE2          = "poe2"
)

var Bf2 = title.GameTitle{
	ProtocolScheme: "bf2",
	GameLabel:      "Battlefield 2",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs: []string{
			"+menu", "1",
			"+restart", "1",
		},
		ExecutableName:    "BF2.exe",
		CloseBeforeLaunch: true,
	},
	CmdBuilder: bf2CmdBuilder,
}

var bf2CmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	profileCon, err := GetDefaultUserProfileCon(bf2ProfileFolder)
	if err != nil {
		return nil, err
	}
	playerName, encryptedPassword, err := GetEncryptedProfileConLogin(profileCon)
	if err != nil {
		return nil, err
	}
	password, err := DecryptProfileConPassword(encryptedPassword)
	if err != nil {
		return nil, err
	}
	args := []string{
		"+joinServer", host,
		"+port", port,
		"+playerName", playerName,
		"+playerPassword", password,
	}

	query := u.Query()
	if query != nil && query.Has(UrlQueryKeyMod) {
		mod, err := getValidMod(installPath, bf2ModBasePath, query.Get(UrlQueryKeyMod), bf2ModSpecialForces, bf2ModAIX2, bf2ModPirates, bf2ModPoE2)
		if err != nil {
			return nil, err
		}

		args = append(args,
			"+modPath", fmt.Sprintf("mods/%s", mod),
			"+ignoreAsserts", "1",
		)
	}

	return args, nil
}

var Bf2SF = title.GameTitle{
	ProtocolScheme: "bf2sf",
	GameLabel:      "Battlefield 2: Special Forces",
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2 Special Forces",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs: []string{
			"+menu", "1",
			"+restart", "1",
			"+modPath", "mods/xpack",
			"+ignoreAsserts", "1",
		},
		ExecutableName:    "BF2.exe",
		CloseBeforeLaunch: true,
	},
	CmdBuilder: bf2CmdBuilder,
}
