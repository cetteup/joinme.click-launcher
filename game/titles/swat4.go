package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
)

var Swat4Config = game.LauncherConfig{
	ProtocolScheme:    "swat4",
	GameLabel:         "SWAT 4",
	ExecutablePath:    "Content\\System\\Swat4.exe",
	StartIn:           game.BinaryDir,
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4",
	RegistryValueName: "InstallPath",
}

var Swat4CmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	return []string{fmt.Sprintf("%s:%s", ip, port)}, nil
}

var Swat4XConfig = game.LauncherConfig{
	ProtocolScheme:    "swat4x",
	GameLabel:         "SWAT 4: The Stetchkov Syndicate",
	ExecutablePath:    "ContentExpansion\\System\\Swat4X.exe",
	StartIn:           game.BinaryDir,
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4 - THE STETCHKOV SYNDICATE",
	RegistryValueName: "InstallPath",
}
