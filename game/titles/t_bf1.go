package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/platform"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

var Bf1 = title.GameTitle{
	ProtocolScheme: "bf1",
	GameLabel:      "Battlefield 1",
	PlatformClient: &platform.OriginClient,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1",
			RegistryValueName: "Install Dir",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "bf1.exe",
	},
	CmdBuilder: func(scheme string, host string, port string, u *url.URL) ([]string, error) {
		offerIDs := []string{"1026023"}
		args := append(frostbite3DefaultArgs, "-gameId", host)

		originURL := buildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
}
