package titles

import "github.com/cetteup/joinme.click-launcher/game"

var Bf1942Config = game.LauncherConfig{
	ProtocolScheme:    "bf1942",
	GameLabel:         "Battlefield 1942",
	ExecutablePath:    "BF1942.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
	RegistryValueName: "GAMEDIR",
}

var Bf1942CmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	args := []string{
		"+restart", "1",
		"+joinServer", ip,
		"+port", port,
	}

	return args, nil
}
