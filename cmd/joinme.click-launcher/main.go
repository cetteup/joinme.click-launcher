//go:generate goversioninfo

package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	filerepo "github.com/cetteup/filerepo/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/router"
	"github.com/cetteup/joinme.click-launcher/internal/titles"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/registry_repository"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	err := internal.LoadConfig()
	if err != nil {
		log.Err(err).Msg("Failed to load configuration from file, continuing with defaults")
	}

	registryRepository := registry_repository.New()
	fileRepository := filerepo.New()

	gameFinder := software_finder.New(registryRepository, fileRepository)
	gameLauncher := game_launcher.New(fileRepository)
	gameRouter = router.New(registryRepository, gameFinder, gameLauncher)
	gameRouter.AddTitle(
		titles.Bf1942,
		titles.BfVietnam,
		titles.Bf2,
		titles.Bf4,
		titles.Bf1,
		titles.Cod,
		titles.CodUO,
		titles.Cod2,
		titles.Cod4,
		titles.CodWaw,
		titles.FearSec2,
		titles.Paraworld,
		titles.Swat4,
		titles.Swat4X,
		titles.Unreal,
		titles.UT,
		titles.UT2003,
		titles.UT2004,
		titles.Vietcong,
	)
}

var (
	gameRouter *router.GameRouter
)

func main() {
	var quietLaunch bool
	var debug bool
	flag.BoolVar(&quietLaunch, "quiet", false, "do not leave the window open any longer than required")
	flag.BoolVar(&debug, "debug", false, "set log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug || internal.Config.DebugLogging {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	args := flag.Args()
	if len(args) == 0 {
		results := gameRouter.RegisterHandlers()
		sort.Slice(results, func(i, j int) bool {
			return results[i].Title.ProtocolScheme < results[j].Title.ProtocolScheme
		})
		for _, result := range results {
			var message string
			if result.Error != nil {
				message = "handler registration failed"
			} else if !result.GameInstalled {
				message = "not installed"
			} else if result.GameInstalled && result.Title.RequiresPlatformClient() && !result.PlatformClientInstalled {
				message = fmt.Sprintf("installed, but required platform client is missing (%s)", result.Title.PlatformClient.Platform)
			} else if result.PreviouslyRegistered {
				message = "launcher already registered"
			} else {
				message = "launcher registered successfully"
			}
			log.Info().
				Err(result.Error).
				Str("game", result.Title.Name).
				Str("result", message).
				Msg("Checked status for")
		}
	} else if len(args) == 1 {
		title, err := gameRouter.RunURL(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Str("game", title.String()).
				Str("url", args[0]).
				Msg("Game could not be launched")
		} else {
			log.Info().
				Str("game", title.String()).
				Str("url", args[0]).
				Msg("Game launched")
		}
	}

	// Leave window open for a bit unless disabled via arg or config
	if !quietLaunch && !internal.Config.QuietLaunch {
		log.Info().Msg("Window will close in 15 seconds")
		time.Sleep(15 * time.Second)
	}
}
