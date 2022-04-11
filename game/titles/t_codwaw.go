package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var CodWawConfig = launcher.Config{
	ProtocolScheme:    "codwaw",
	GameLabel:         "Call of Duty: World at War",
	ExecutablePath:    "CoDWaWmp.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty WAW",
	RegistryValueName: "InstallPath",
}
