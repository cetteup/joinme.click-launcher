package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var Swat4Config = launcher.Config{
	ProtocolScheme:    "swat4",
	GameLabel:         "SWAT 4",
	ExecutablePath:    "Content\\System\\Swat4.exe",
	StartIn:           launcher.BinaryDir,
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4",
	RegistryValueName: "InstallPath",
}

var Swat4CmdBuilder launcher.CommandBuilder = func(config launcher.Config, ip string, port string) ([]string, error) {
	return []string{fmt.Sprintf("%s:%s", ip, port)}, nil
}

var Swat4XConfig = launcher.Config{
	ProtocolScheme:    "swat4x",
	GameLabel:         "SWAT 4: The Stetchkov Syndicate",
	ExecutablePath:    "ContentExpansion\\System\\Swat4X.exe",
	StartIn:           launcher.BinaryDir,
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Sierra\\SWAT 4 - THE STETCHKOV SYNDICATE",
	RegistryValueName: "InstallPath",
}
