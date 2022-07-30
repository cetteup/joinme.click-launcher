//go:generate goversioninfo

package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/router"
	"github.com/cetteup/joinme.click-launcher/internal/titles"
	"github.com/cetteup/joinme.click-launcher/pkg/registry_repository"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	err := internal.LoadConfig()
	if err != nil {
		log.Err(err).Msg("Failed to load configuration from file, continuing with defaults")
	}

	registryRepository := registry_repository.NewRegistryRepository()

	gameFinder := software_finder.NewSoftwareFinder(registryRepository)
	gameRouter = router.NewGameRouter(registryRepository, gameFinder)
	gameRouter.AddTitle(titles.Bf1942)
	gameRouter.AddTitle(titles.Bf1942RoadToRome)
	gameRouter.AddTitle(titles.Bf1942SecretWeaponsOfWW2)
	gameRouter.AddTitle(titles.BfVietnam)
	gameRouter.AddTitle(titles.Bf2)
	gameRouter.AddTitle(titles.Bf2SF)
	gameRouter.AddTitle(titles.Bf4)
	gameRouter.AddTitle(titles.Bf1)
	gameRouter.AddTitle(titles.Cod)
	gameRouter.AddTitle(titles.CodUO)
	gameRouter.AddTitle(titles.Cod2)
	gameRouter.AddTitle(titles.Cod4)
	gameRouter.AddTitle(titles.CodWaw)
	gameRouter.AddTitle(titles.FearSec2)
	gameRouter.AddTitle(titles.Paraworld)
	gameRouter.AddTitle(titles.Swat4)
	gameRouter.AddTitle(titles.Swat4X)
	gameRouter.AddTitle(titles.Vietcong)
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
				Str("game", result.Title.GameLabel).
				Str("result", message).
				Msg("Checked status for")
		}
	} else if len(args) == 1 {
		err := gameRouter.RunURL(args[0])
		if err != nil {
			log.Error().
				Err(err).
				Str("url", args[0]).
				Msg("Game could not be launched")
		} else {
			log.Info().
				Str("url", args[0]).
				Msg("Successfully launched game")
		}
	}

	// Leave window open for a bit unless disabled via arg or config
	if !quietLaunch && !internal.Config.QuietLaunch {
		log.Info().Msg("Window will close in 15 seconds")
		time.Sleep(15 * time.Second)
	}
}
