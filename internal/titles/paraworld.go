package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	paraworldModPathTemplate = "Data\\%s\\UI\\All_def.txt"
	paraworldModBoosterPack  = "BoosterPack1"
	paraworldModMirage       = "MIRAGE"
)

var Paraworld = domain.GameTitle{
	ProtocolScheme: "paraworld",
	GameLabel:      "ParaWorld",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName:    "Paraworld.exe",
		ExecutablePath:    "bin",
		CloseBeforeLaunch: true,
		AdditionalProcessNames: map[string]bool{
			"PWClient.exe": true,
			"PWServer.exe": true,
		},
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder: func(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		var args []string
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, "-autoconnect", fmt.Sprintf("%s:%s", u.Hostname(), u.Port()))
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			mod, err := internal.GetValidModFromQuery(
				query,
				config.InstallPath,
				paraworldModPathTemplate,
				software_finder.PathTypeFile,
				paraworldModBoosterPack, paraworldModMirage,
			)
			if err != nil {
				return nil, err
			}

			args = append(args, "-enable", mod)
		}
		return args, nil
	},
}
