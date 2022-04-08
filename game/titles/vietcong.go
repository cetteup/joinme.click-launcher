package titles

import (
	"github.com/cetteup/joinme.click-launcher/game"
)

var VietcongConfig = game.LauncherConfig{
	ProtocolScheme:    "vietcong",
	GameLabel:         "Vietcong",
	ExecutablePath:    "vietcong.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Pterodon\\Vietcong",
	RegistryValueName: "InstallDir",
}

var VietcongCmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	return []string{"-ip", ip, "-port", port}, nil
}
