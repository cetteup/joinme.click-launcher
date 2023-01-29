package titles

import (
	"fmt"

	"github.com/cetteup/joinme.click-launcher/internal/domain"
	localinternal "github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bf1942ModPathTemplate       = "Mods\\%s\\lexiconAll.dat"
	bf1942ModRoadToRome         = "Xpack1"
	bf1942ModSecretWeaponsOfWW2 = "Xpack2"
	bf1942Mod1918               = "bf1918"
	bf1942ModDCFinal            = "DC_Final"
	bf1942ModDesertCombat       = "DesertCombat"
	bf1942ModPirates            = "Pirates"
)

var Bf1942 = domain.GameTitle{
	Name:           "Battlefield 1942",
	ProtocolScheme: "bf1942",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942",
			RegistryValueName: "GAMEDIR",
		},
	},
	Mods: []domain.GameMod{
		domain.MakeMod(
			"The Road to Rome",
			bf1942ModRoadToRome,
			[]software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack1",
					RegistryValueName: "GAMEDIR",
				},
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942ModRoadToRome),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Secret Weapons of WWII",
			bf1942ModSecretWeaponsOfWW2,
			[]software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\WOW6432Node\\EA Games\\Battlefield 1942 Xpack2",
					RegistryValueName: "GAMEDIR",
				},
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942ModSecretWeaponsOfWW2),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Battlefield 1918",
			bf1942Mod1918,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942Mod1918),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Desert Combat Final",
			bf1942ModDCFinal,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942ModDCFinal),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Desert Combat (0.7)",
			bf1942ModDesertCombat,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942ModDesertCombat),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Pirates",
			bf1942ModPirates,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf1942ModPathTemplate, bf1942ModPirates),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs:    []string{"+restart", "1"},
		ExecutableName: "BF1942.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     localinternal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: localinternal.IPPortURLValidator{},
	CmdBuilder:   localinternal.RefractorV1CmdBuilder{},
	HookHandlers: []game_launcher.HookHandler{
		localinternal.MakeKillProcessHookHandler(true),
	},
}
