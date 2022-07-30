package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/platforms"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Bf1 = domain.GameTitle{
	ProtocolScheme: "bf1",
	GameLabel:      "Battlefield 1",
	PlatformClient: &platforms.OriginClient,
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1",
			RegistryValueName: "Install Dir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "bf1.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.Frostbite3GameIdURLValidator,
	CmdBuilder: func(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		var args []string
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, internal.Frostbite3DefaultArgs...)
			args = append(args, "-gameId", u.Hostname())
		}

		offerIDs := []string{"1026023"}
		originURL := internal.BuildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
}