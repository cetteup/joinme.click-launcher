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
	bfVietnamModPathTemplate  = "Mods\\%s\\lexiconAll.dat"
	bfVietnamModBattlegroup42 = "Battlegroup42"
)

var BfVietnam = domain.GameTitle{
	Name:           "Battlefield Vietnam",
	ProtocolScheme: "bfvietnam",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield Vietnam",
			RegistryValueName: "GAMEDIR",
		},
	},
	Mods: []domain.GameMod{
		domain.MakeMod(
			"Battlegroup 42",
			bfVietnamModBattlegroup42,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bfVietnamModPathTemplate, bfVietnamModBattlegroup42),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs:    []string{"+restart", "1"},
		ExecutableName: "BfVietnam.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     localinternal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: localinternal.IPPortURLValidator,
	CmdBuilder: func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		args := config.DefaultArgs
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, "+joinServer", u.Hostname(), "+port", u.Port())
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			args = append(args, "+game", internal.GetModFromQuery(query))
		}

		return args, nil
	},
	HookHandlers: map[string]game_launcher.HookHandler{
		localinternal.HookKillProcess: localinternal.KillProcessHookHandler(true),
	},
}
