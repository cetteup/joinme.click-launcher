package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var VietcongConfig = launcher.Config{
	ProtocolScheme:    "vietcong",
	GameLabel:         "Vietcong",
	ExecutablePath:    "vietcong.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Pterodon\\Vietcong",
	RegistryValueName: "InstallDir",
}

var VietcongCmdBuilder launcher.CommandBuilder = func(config launcher.Config, ip string, port string) ([]string, error) {
	return []string{"-ip", ip, "-port", port}, nil
}
