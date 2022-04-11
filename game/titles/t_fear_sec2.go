package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var FearSec2Config = launcher.Config{
	ProtocolScheme:    "fearsec2",
	GameLabel:         "F.E.A.R. Combat (SEC2)",
	ExecutablePath:    "FEARMP.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\FEAR-Community.org\\FEAR Combat (SEC2)",
	RegistryValueName: "Path",
}

var FearSec2CmdBuilder launcher.CommandBuilder = func(config launcher.Config, ip string, port string) ([]string, error) {
	return []string{"+join", fmt.Sprintf("%s:%s", ip, port)}, nil
}
