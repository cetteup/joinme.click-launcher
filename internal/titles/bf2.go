package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/conman/pkg/game/bf2"
	"github.com/cetteup/conman/pkg/handler"
	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	localinternal "github.com/cetteup/joinme.click-launcher/internal/titles/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

const (
	bf2ModPathTemplate  = "mods\\%s\\Common_client.zip"
	bf2ModSpecialForces = "xpack"
	bf2ModAIX2          = "AIX2"
	bf2ModArcticWarfare = "Arctic_Warfare"
	bf2ModPirates       = "bfp2"
	bf2ModPoE2          = "poe2"
)

var Bf2 = domain.GameTitle{
	Name:           "Battlefield 2",
	ProtocolScheme: "bf2",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryKey:       software_finder.RegistryKeyLocalMachine,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
			RegistryValueName: "InstallDir",
		},
	},
	Mods: []domain.GameMod{
		domain.MakeMod(
			"Special Forces",
			bf2ModSpecialForces,
			[]software_finder.Config{
				{
					ForType:           software_finder.RegistryFinder,
					RegistryKey:       software_finder.RegistryKeyLocalMachine,
					RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2 Special Forces",
					RegistryValueName: "InstallDir",
				},
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf2ModPathTemplate, bf2ModSpecialForces),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Allied Intent Xtended",
			bf2ModAIX2,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf2ModPathTemplate, bf2ModAIX2),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Pirates (Yarr2)",
			bf2ModPirates,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf2ModPathTemplate, bf2ModPirates),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Point of Existence 2",
			bf2ModPoE2,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf2ModPathTemplate, bf2ModPoE2),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
		domain.MakeMod(
			"Arctic Warfare",
			bf2ModArcticWarfare,
			[]software_finder.Config{
				{
					ForType:     software_finder.PathFinder,
					InstallPath: fmt.Sprintf(bf2ModPathTemplate, bf2ModArcticWarfare),
					PathType:    software_finder.PathTypeFile,
				},
			},
		),
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs: []string{
			"+menu", "1",
			"+restart", "1",
		},
		ExecutableName:    "BF2.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: localinternal.IPPortURLValidator,
	CmdBuilder: func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		configHandler := handler.New(fr)
		profileCon, err := bf2.GetDefaultProfileProfileCon(configHandler)
		if err != nil {
			return nil, err
		}

		playerName, encryptedPassword, err := bf2.GetEncryptedLogin(profileCon)
		if err != nil {
			return nil, fmt.Errorf("failed to extract login details from profile.con: %s", err)
		}

		password, err := bf2.DecryptProfileConPassword(encryptedPassword)
		if err != nil {
			return nil, err
		}

		args := append(config.DefaultArgs, "+playerName", playerName, "+playerPassword", password)
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, "+joinServer", u.Hostname(), "+port", u.Port())
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			args = append(args,
				"+modPath", fmt.Sprintf("mods/%s", internal.GetModFromQuery(query)),
				"+ignoreAsserts", "1",
			)
		}

		return args, nil
	},
}
