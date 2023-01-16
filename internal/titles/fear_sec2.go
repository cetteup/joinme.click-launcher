package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var FearSec2 = domain.GameTitle{
	Name:           "F.E.A.R. Combat (SEC2)",
	ProtocolScheme: "fearsec2",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\FEAR-Community.org\\FEAR Combat (SEC2)",
			RegistryValueName: "Path",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "FEARMP.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder: func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			return append(config.DefaultArgs, "+join", fmt.Sprintf("%s:%s", u.Hostname(), u.Port())), nil
		}
		return nil, nil
	},
}
