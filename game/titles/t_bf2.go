package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

const (
	ProfileFolder = "Battlefield 2"
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
		ExecutablePath: "BF2.exe",
	},
	CmdBuilder: bf2CmdBuilder,
}

var bf2CmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	profileCon, err := GetDefaultUserProfileCon(ProfileFolder)
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
		ExecutablePath: "BF2.exe",
	},
	CmdBuilder: bf2CmdBuilder,
}
