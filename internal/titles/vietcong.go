package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var Vietcong = domain.GameTitle{
	Name:           "Vietcong",
	ProtocolScheme: "vietcong",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Pterodon\\Vietcong",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "vietcong.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator{},
	CmdBuilder:   vietcongCmdBuilder{},
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
	},
}

type vietcongCmdBuilder struct{}

func (b vietcongCmdBuilder) GetArgs(_ game_launcher.FileRepository, u *url.URL, launchType game_launcher.LaunchType) ([]string, error) {
	args := make([]string, 0, 4)
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args, "-ip", u.Hostname(), "-port", u.Port())
	}

	return args, nil
}
