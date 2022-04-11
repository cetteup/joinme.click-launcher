package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var ParaworldConfig = launcher.Config{
	ProtocolScheme:    "paraworld",
	GameLabel:         "ParaWorld",
	ExecutablePath:    "bin\\Paraworld.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
	RegistryValueName: "InstallDir",
}

var ParaworldCmdBuilder launcher.CommandBuilder = func(config launcher.Config, ip string, port string) ([]string, error) {
	return []string{"-autoconnect", fmt.Sprintf("%s:%s", ip, port)}, nil
}
