package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	localinternal "github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	paraworldModPathTemplate = "Data\\%s\\UI\\All_def.txt"
	paraworldModBoosterPack  = "BoosterPack1"
	paraworldModMirage       = "MIRAGE"
)

var Paraworld = domain.GameTitle{
	Name:           "ParaWorld",
	ProtocolScheme: "paraworld",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Sunflowers\\ParaWorld",
			RegistryValueName: "InstallDir",
		},
	},
	Mods: []domain.GameMod{
		domain.MakeMod(
			"Booster pack",
			paraworldModBoosterPack,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(paraworldModPathTemplate, paraworldModBoosterPack),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Mirage",
			paraworldModMirage,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(paraworldModPathTemplate, paraworldModMirage),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "Paraworld.exe",
		ExecutablePath: "bin",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     localinternal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: localinternal.IPPortURLValidator,
	CmdBuilder:   paraworldCmdBuilder{},
	HookHandlers: []game_launcher.HookHandler{
		localinternal.MakeKillProcessHookHandler(true, "PWClient.exe", "PWServer.exe"),
	},
}

type paraworldCmdBuilder struct{}

func (b paraworldCmdBuilder) GetArgs(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	args := config.DefaultArgs
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args, "-autoconnect", fmt.Sprintf("%s:%s", u.Hostname(), u.Port()))
	}

	query := u.Query()
	if internal.QueryHasMod(query) {
		args = append(args, "-enable", internal.GetModFromQuery(query))
	}
	return args, nil
}
