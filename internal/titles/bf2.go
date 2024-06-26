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

	bf2HookPurgeServerHistory = "purge-server-history"
	bf2HookPurgeShaderCache   = "purge-shader-cache"
	bf2HookPurgeLogoCache     = "purge-logo-cache"
	bf2HookSetDefaultProfile  = "set-default-profile"
	hookArgProfile            = "profile"
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
		ExecutableName: "BF2.exe",
		HookConfigs: []game_launcher.HookConfig{
			{
				Handler:     localinternal.HookKillProcess,
				When:        game_launcher.HookWhenPreLaunch,
				ExitOnError: true,
			},
		},
	},
	URLValidator: localinternal.IPPortURLValidator{},
	CmdBuilder:   bf2CmdBuilder{},
	HookHandlers: []game_launcher.HookHandler{
		localinternal.MakeKillProcessHookHandler(true),
		bf2SetDefaultProfileHookHandler{},
		bf2PurgeServerHistoryHookHandler{},
		bf2PurgeShaderCacheHookHandler{},
		bf2PurgeLogoCacheHookHandler{},
	},
}

type bf2CmdBuilder struct{}

func (b bf2CmdBuilder) GetArgs(fr game_launcher.FileRepository, u *url.URL, launchType game_launcher.LaunchType) ([]string, error) {
	configHandler := handler.New(fr)
	profileCon, err := bf2.GetDefaultProfileProfileCon(configHandler)
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, 12)
	// Only multiplayer profiles contain an email address
	if profileCon.HasKey(bf2.ProfileConKeyEmail) {
		playerName, encryptedPassword, err2 := bf2.GetEncryptedLogin(profileCon)
		if err2 != nil {
			return nil, fmt.Errorf("failed to extract login details from profile.con: %s", err)
		}

		password, err2 := bf2.DecryptProfileConPassword(encryptedPassword)
		if err2 != nil {
			return nil, fmt.Errorf("failed to decrypt player password: %s", err)
		}

		args = append(args, "+playerName", playerName, "+playerPassword", password)
	} else {
		// Singleplayer profiles always have an empty GamespyNick, so use the "normal" nick instead
		playerName, err2 := profileCon.GetValue(bf2.ProfileConKeyNick)
		if err2 != nil {
			return nil, fmt.Errorf("failed to extract player name from profile.con: %s", err2)
		}

		args = append(args, "+playerName", playerName.String())
	}

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
}

type bf2SetDefaultProfileHookHandler struct{}

func (h bf2SetDefaultProfileHookHandler) Run(fr game_launcher.FileRepository, _ *url.URL, _ game_launcher.Config, _ game_launcher.LaunchType, args map[string]string) error {
	profileKey, ok := args[hookArgProfile]
	if !ok {
		return fmt.Errorf("required argument %s for hook %s is missing", hookArgProfile, h.String())
	}

	configHandler := handler.New(fr)
	globalCon, err := configHandler.ReadGlobalConfig(handler.GameBf2)
	if err != nil {
		return err
	}

	bf2.SetDefaultProfile(globalCon, profileKey)

	return configHandler.WriteConfigFile(globalCon)
}

func (h bf2SetDefaultProfileHookHandler) String() string {
	return bf2HookSetDefaultProfile
}

type bf2PurgeServerHistoryHookHandler struct{}

func (h bf2PurgeServerHistoryHookHandler) Run(fr game_launcher.FileRepository, _ *url.URL, _ game_launcher.Config, _ game_launcher.LaunchType, args map[string]string) error {
	configHandler := handler.New(fr)
	profileKey, ok := args[hookArgProfile]
	if !ok {
		// Use default profile if none has been configured
		var err error
		profileKey, err = bf2.GetDefaultProfileKey(configHandler)
		if err != nil {
			return err
		}
	}

	generalCon, err := bf2.ReadProfileConfigFile(configHandler, profileKey, bf2.ProfileConfigFileGeneralCon)
	if err != nil {
		return err
	}

	bf2.PurgeServerHistory(generalCon)

	return configHandler.WriteConfigFile(generalCon)
}

func (h bf2PurgeServerHistoryHookHandler) String() string {
	return bf2HookPurgeServerHistory
}

type bf2PurgeShaderCacheHookHandler struct{}

func (h bf2PurgeShaderCacheHookHandler) Run(fr game_launcher.FileRepository, _ *url.URL, _ game_launcher.Config, _ game_launcher.LaunchType, _ map[string]string) error {
	configHandler := handler.New(fr)
	return configHandler.PurgeShaderCache(handler.GameBf2)
}

func (h bf2PurgeShaderCacheHookHandler) String() string {
	return bf2HookPurgeShaderCache
}

type bf2PurgeLogoCacheHookHandler struct{}

func (h bf2PurgeLogoCacheHookHandler) Run(fr game_launcher.FileRepository, _ *url.URL, _ game_launcher.Config, _ game_launcher.LaunchType, _ map[string]string) error {
	configHandler := handler.New(fr)
	return configHandler.PurgeLogoCache(handler.GameBf2)
}

func (h bf2PurgeLogoCacheHookHandler) String() string {
	return bf2HookPurgeLogoCache
}
