package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bfVietnamModPathTemplate  = "Mods\\%s\\lexiconAll.dat"
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
	CmdBuilder: func(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		var args []string
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, "+joinServer", u.Hostname(), "+port", u.Port())
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			mod, err := internal.GetValidModFromQuery(
				query,
				config.InstallPath,
				bfVietnamModPathTemplate,
				software_finder.PathTypeFile,
				bfVietnamModBattlegroup42,
			)
			if err != nil {
				return nil, err
			}

			args = append(args, "+game", mod)
		}

		return args, nil
	},
}
