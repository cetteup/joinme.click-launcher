package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
)

const (
	bfVietnamModBasePath      = "Mods"
	bfVietnamModBattlegroup42 = "Battlegroup42"
)

var BfVietnam = title.GameTitle{
	ProtocolScheme: "bfvietnam",
	GameLabel:      "Battlefield Vietnam",
	RequiresPort:   true,
	FinderConfigs: []finder.Config{
		{
			ForType:           finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield Vietnam",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: launcher.Config{
		DefaultArgs:       []string{"+restart", "1"},
		ExecutableName:    "BfVietnam.exe",
		CloseBeforeLaunch: true,
	},
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		args := []string{
			"+joinServer", host,
			"+port", port,
		}

		query := u.Query()
		if query != nil && query.Has(urlQueryKeyMod) {
			mod, err := getValidMod(installPath, bfVietnamModBasePath, query.Get(urlQueryKeyMod), bfVietnamModBattlegroup42)
			if err != nil {
				return nil, err
			}

			args = append(args, "+game", mod)
		}

		return args, nil
	},
}
