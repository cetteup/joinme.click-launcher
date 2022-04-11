package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var Cod4Config = launcher.Config{
	ProtocolScheme:    "cod4",
	GameLabel:         "Call of Duty 4: Modern Warfare",
	ExecutablePath:    "iw3mp.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 4",
	RegistryValueName: "InstallPath",
}
