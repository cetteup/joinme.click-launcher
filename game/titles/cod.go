package titles

import "github.com/cetteup/joinme.click-launcher/game"

var CodConfig = game.LauncherConfig{
	ProtocolScheme:    "cod",
	GameLabel:         "Call of Duty",
	ExecutablePath:    "CoDMP.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty",
	RegistryValueName: "InstallPath",
}
