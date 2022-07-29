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
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		offerIDs := []string{"1026023"}
		args := append(internal.Frostbite3DefaultArgs, "-gameId", host)

		originURL := internal.BuildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
}
