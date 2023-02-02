package titles

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

var CodWaw = domain.GameTitle{
	Name:           "Call of Duty: World at War",
	ProtocolScheme: "codwaw",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Activision\\Call of Duty WAW",
			RegistryValueName: "InstallPath",
		},
	},
	LauncherConfig: game_launcher.Config{
		ExecutableName: "CoDWaWmp.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     internal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
			{
				Handler:     internal.HookDeleteFile,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: false,
			},
		},
	},
	URLValidator: internal.IPPortURLValidator{},
	CmdBuilder:   internal.MakeSimpleCmdBuilder(internal.PlusConnectPrefix),
	HookHandlers: []game_launcher.HookHandler{
		internal.MakeKillProcessHookHandler(true),
		internal.MakeDeleteFileHookHandler(codWawRunningFilePathsBuilder),
	},
}

var codWawRunningFilePathsBuilder = func(config game_launcher.Config) ([]string, error) {
	name := fmt.Sprintf("__%s", strings.TrimSuffix(config.ExecutableName, filepath.Ext(config.ExecutableName)))
	appData, err := internal.GetLocalAppDataPath()
	if err != nil {
		return nil, err
	}
	return []string{filepath.Join(appData, "Activision", "CoDWaW", name)}, nil
}
