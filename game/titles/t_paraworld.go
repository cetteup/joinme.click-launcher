package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

const (
	paraworldModBasePath    = "Data"
	paraworldModBoosterPack = "BoosterPack1"
	paraworldModMirage      = "MIRAGE"
)

var Paraworld = title.GameTitle{
	ProtocolScheme: "paraworld",
	GameLabel:      "ParaWorld",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: launcher.Config{
		ExecutableName:    "Paraworld.exe",
		ExecutablePath:    "bin",
		CloseBeforeLaunch: true,
		AdditionalProcessNames: map[string]bool{
			"PWClient.exe": true,
			"PWServer.exe": true,
		},
	},
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		args := []string{"-autoconnect", fmt.Sprintf("%s:%s", host, port)}
		query := u.Query()
		if query != nil && query.Has(UrlQueryKeyMod) {
			mod, err := getValidMod(installPath, paraworldModBasePath, query.Get(UrlQueryKeyMod), paraworldModBoosterPack, paraworldModMirage)
			if err != nil {
				return nil, err
			}

			args = append(args, "-enable", mod)
		}
		return args, nil
	},
}
