package titles

import (
	"net/url"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/platforms"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bf1Exe = "bf1.exe"
)

var Bf1 = domain.GameTitle{
	Name:           "Battlefield 1",
	ProtocolScheme: "bf1",
	PlatformClient: &platforms.OriginClient,
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1",
			RegistryValueName: "Install Dir",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: bf1Exe,
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: internal.Frostbite3GameIdURLValidator,
	CmdBuilder: func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		args := config.DefaultArgs
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, internal.Frostbite3DefaultArgs...)
			args = append(args, "-gameId", u.Hostname())
		}

		offerIDs := []string{"1026023"}
		originURL := internal.BuildOriginURL(offerIDs, args)
		return []string{originURL}, nil
	},
	HookHandlers: map[string]game_launcher.HookHandler{
		internal.HookKillProcess: internal.KillProcessHookHandler(false, bf1Exe), // Launcher config executable name will be "Origin.exe", which we don't want to kill
	},
}
