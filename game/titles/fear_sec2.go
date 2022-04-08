package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
)

var FearSec2Config = game.LauncherConfig{
	ProtocolScheme:    "fearsec2",
	GameLabel:         "F.E.A.R. Combat (SEC2)",
	ExecutablePath:    "FEARMP.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\FEAR-Community.org\\FEAR Combat (SEC2)",
	RegistryValueName: "Path",
}

var FearSec2CmdBuilder game.CommandBuilder = func(config game.LauncherConfig, ip string, port string) ([]string, error) {
	return []string{"+join", fmt.Sprintf("%s:%s", ip, port)}, nil
}
