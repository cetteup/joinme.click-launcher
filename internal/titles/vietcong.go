package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Vietcong = domain.GameTitle{
	ProtocolScheme: "vietcong",
	GameLabel:      "Vietcong",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Pterodon\\Vietcong",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "vietcong.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder: func(u *url.URL, config game_launcher.Config) ([]string, error) {
		return []string{"-ip", u.Hostname(), "-port", u.Port()}, nil
	},
}
