package titles

import "github.com/cetteup/joinme.click-launcher/game"

var CodConfig = game.LauncherConfig{
	ProtocolScheme:    "cod",
	GameLabel:         "Call of Duty",
	ExecutablePath:    "CoDMP.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty",
	RegistryValueName: "InstallPath",
}

var CodUOConfig = game.LauncherConfig{
	ProtocolScheme:    "coduo",
	GameLabel:         "Call of Duty: United Offensive",
	ExecutablePath:    "CoDUOMP.exe",
	RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty United Offensive",
	RegistryValueName: "InstallPath",
}
