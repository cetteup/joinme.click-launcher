package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var Bf1942Config = launcher.Config{
	ProtocolScheme:    "bf1942",
	GameLabel:         "Battlefield 1942",
	ExecutablePath:    "BF1942.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
	RegistryValueName: "GAMEDIR",
}

var Bf1942CmdBuilder launcher.CommandBuilder = func(config launcher.Config, ip string, port string) ([]string, error) {
	args := []string{
		"+restart", "1",
		"+joinServer", ip,
		"+port", port,
	}

	return args, nil
}
