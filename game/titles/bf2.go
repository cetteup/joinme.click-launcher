package titles

import (
	"github.com/cetteup/joinme.click-launcher/game"
)

var Bf2Config = game.LauncherConfig{
	ProtocolScheme:    "bf2",
	GameLabel:         "Battlefield 2",
	ExecutablePath:    "BF2.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
	RegistryValueName: "InstallDir",
}

var Bf2CmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	profileCon, err := game.GetDefaultUserProfileCon(config.GameLabel)
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
		"+restart", "1",
		"+joinServer", ip,
		"+port", port,
		"+playerName", playerName,
		"+playerPassword", password,
	}

	return args, nil
}
