package titles

import "github.com/cetteup/joinme.click-launcher/game"

var BfVietnamConfig = game.LauncherConfig{
	ProtocolScheme:    "bfvietnam",
	GameLabel:         "Battlefield Vietnam",
	ExecutablePath:    "BfVietnam.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield Vietnam",
	RegistryValueName: "GAMEDIR",
}

var BfVietnamCmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	args := []string{
		"+restart", "1",
		"+joinServer", ip,
		"+port", port,
	}

	return args, nil
}
