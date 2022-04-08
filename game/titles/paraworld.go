package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
)

var ParaworldConfig = game.LauncherConfig{
	ProtocolScheme:    "paraworld",
	GameLabel:         "ParaWorld",
	ExecutablePath:    "bin\\Paraworld.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
	RegistryValueName: "InstallDir",
}

var ParaworldCmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	return []string{"-autoconnect", fmt.Sprintf("%s:%s", ip, port)}, nil
}
