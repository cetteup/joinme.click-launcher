package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bfVietnamModBasePath      = "Mods"
	bfVietnamModBattlegroup42 = "Battlegroup42"
)

var BfVietnam = domain.GameTitle{
	ProtocolScheme: "bfvietnam",
	GameLabel:      "Battlefield Vietnam",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield Vietnam",
			RegistryValueName: "GAMEDIR",
		},
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs:       []string{"+restart", "1"},
		ExecutableName:    "BfVietnam.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder: func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
		args := []string{
			"+joinServer", host,
			"+port", port,
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			mod, err := internal.GetValidModFromQuery(query, installPath, bfVietnamModBasePath, bfVietnamModBattlegroup42)
			if err != nil {
				return nil, err
			}

			args = append(args, "+game", mod)
		}

		return args, nil
	},
}
