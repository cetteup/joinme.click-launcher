package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
)

var CodWawConfig = game.LauncherConfig{
	ProtocolScheme:    "codwaw",
	GameLabel:         "Call of Duty: World at War",
	ExecutablePath:    "CoDWaWmp.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty WAW",
	RegistryValueName: "InstallPath",
}

var CodWawCmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", ip, port)}, nil
}
