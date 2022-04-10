package titles

import (
	"github.com/cetteup/joinme.click-launcher/game"
)

const (
	ProfileFolder = "Battlefield 2"
)

var Bf2Config = game.LauncherConfig{
	ProtocolScheme:    "bf2",
	GameLabel:         "Battlefield 2",
	ExecutablePath:    "BF2.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
	RegistryValueName: "InstallDir",
}

var Bf2CmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	profileCon, err := game.GetDefaultUserProfileCon(ProfileFolder)
	if err != nil {
		return nil, err
	}
	playerName, encryptedPassword, err := game.GetEncryptedProfileConLogin(profileCon)
	if err != nil {
		return nil, err
	}
	password, err := game.DecryptProfileConPassword(encryptedPassword)
	if err != nil {
		return nil, err
	}
	args := []string{
		"+menu", "1",
		"+fullscreen", "1",
		"+restart", "1",
		"+joinServer", ip,
		"+port", port,
		"+playerName", playerName,
		"+playerPassword", password,
	}

	if config.ProtocolScheme == Bf2SFConfig.ProtocolScheme {
		args = append(args, "+modPath", "mods/xpack", "+ignoreAsserts", "1")
	}

	return args, nil
}

var Bf2SFConfig = game.LauncherConfig{
	ProtocolScheme:    "bf2sf",
	GameLabel:         "Battlefield 2: Special Forces",
	ExecutablePath:    "BF2.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2 Special Forces",
	RegistryValueName: "InstallDir",
}
