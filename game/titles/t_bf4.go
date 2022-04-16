package titles

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/platform"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"net/url"
)

var Bf4 = title.GameTitle{
	ProtocolScheme: "bf4",
	GameLabel:      "Battlefield 4",
	PlatformClient: &platform.OriginClient,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 4",
			RegistryValueName: "Install Dir",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutablePath: "bf4.exe",
	},
	CmdBuilder: func(scheme string, host string, port string, u *url.URL) ([]string, error) {
		offerIDs := []string{"1007968", "1011575", "1011576", "1011577", "1010268", "1010269", "1010270", "1010271", "1010958", "1010959", "1010960", "1010961", "1007077", "1016751", "1016757", "1016754", "1015365", "1015364", "1015363", "1015362"}
		args := append(frostbite3DefaultArgs, "-gameId", host)

		originURL := buildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
}
