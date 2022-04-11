package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var Cod2Config = launcher.Config{
	ProtocolScheme:    "cod2",
	GameLabel:         "Call of Duty 2",
	ExecutablePath:    "CoD2MP_s.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty 2",
	RegistryValueName: "InstallPath",
}
