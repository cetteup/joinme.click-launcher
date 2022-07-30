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
	bf2ProfileFolder    = "Battlefield 2"
	bf2ModPathTemplate  = "mods\\%s\\Common_client.zip"
	bf2ModSpecialForces = "xpack"
	bf2ModAIX2          = "AIX2"
	bf2ModPirates       = "bfp2"
	bf2ModPoE2          = "poe2"
)

var Bf2 = domain.GameTitle{
	ProtocolScheme: "bf2",
	GameLabel:      "Battlefield 2",
	FinderConfigs: []software_finder.Config{
		{
			ForType:           software_finder.RegistryFinder,
			RegistryPath:      "SOFTWARE\\WOW6432Node\\Electronic Arts\\EA Games\\Battlefield 2",
			RegistryValueName: "InstallDir",
		},
	},
	LauncherConfig: game_launcher.Config{
		DefaultArgs: []string{
			"+menu", "1",
			"+restart", "1",
		},
		ExecutableName:    "BF2.exe",
		CloseBeforeLaunch: true,
	},
	URLValidator: internal.IPPortURLValidator,
	CmdBuilder: func(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
		profileCon, err := internal.GetDefaultUserProfileCon(bf2ProfileFolder)
		if err != nil {
			return nil, err
		}

		playerName, encryptedPassword, err := internal.GetEncryptedProfileConLogin(profileCon)
		if err != nil {
			return nil, fmt.Errorf("failed to extract login details from profile.con: %s", err)
		}

		password, err := internal.DecryptProfileConPassword(encryptedPassword)
		if err != nil {
			return nil, err
		}

		args := []string{"+playerName", playerName, "+playerPassword", password}
		if launchType == game_launcher.LaunchTypeLaunchAndJoin {
			args = append(args, "+joinServer", u.Hostname(), "+port", u.Port())
		}

		query := u.Query()
		if internal.QueryHasMod(query) {
			mod, err := internal.GetValidModFromQuery(
				query,
				config.InstallPath,
				bf2ModPathTemplate,
				software_finder.PathTypeFile,
				bf2ModSpecialForces, bf2ModAIX2, bf2ModPirates, bf2ModPoE2,
			)
			if err != nil {
				return nil, err
			}

			args = append(args,
				"+modPath", fmt.Sprintf("mods/%s", mod),
				"+ignoreAsserts", "1",
			)
		}

		return args, nil
	},
}
