package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/platforms"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Bf4 = domain.GameTitle{
	ProtocolScheme: "bf4",
	GameLabel:      "Battlefield 4",
	PlatformClient: &platforms.OriginClient,
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 4",
			RegistryValueName: "Install Dir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "bf4.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.Frostbite3GameIdURLValidator,
	CmdBuilder: func(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		var args []string
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, internal.Frostbite3DefaultArgs...)
			args = append(args, "-gameId", u.Hostname())
		}

		offerIDs := []string{"1007968", "1011575", "1011576", "1011577", "1010268", "1010269", "1010270", "1010271", "1010958", "1010959", "1010960", "1010961", "1007077", "1016751", "1016757", "1016754", "1015365", "1015364", "1015363", "1015362"}
		originURL := internal.BuildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
}
