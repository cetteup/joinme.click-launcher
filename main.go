//go:generate goversioninfo

package main

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/router"
	"github.com/cetteup/joinme.click-launcher/game/titles"
	"github.com/cetteup/joinme.click-launcher/internal"
	"os"
	"sort"
	"time"
)

func init() {
	err := internal.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %s\n", err)
		fmt.Println("Continuing with defaults")
	}

	registryRepository := internal.NewRegistryRepository()

	gameFinder := finder.NewSoftwareFinder(registryRepository)
	gameRouter = router.NewGameRouter(registryRepository, gameFinder)
	gameRouter.AddTitle(titles.Bf1942)
	gameRouter.AddTitle(titles.BfVietnam)
	gameRouter.AddTitle(titles.Bf2)
	gameRouter.AddTitle(titles.Bf2SF)
	gameRouter.AddTitle(titles.Bf4)
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
	args := os.Args[1:]
	if len(args) == 0 {
		results := gameRouter.RegisterHandlers()
		sort.Slice(results, func(i, j int) bool {
			return results[i].Title.ProtocolScheme < results[j].Title.ProtocolScheme
		})
		for _, result := range results {
			var message string
			if result.Error != nil {
				message = fmt.Sprintf("handler registration failed (%s)\n", result.Error)
			} else if !result.GameInstalled {
				message = "not installed"
			} else if result.GameInstalled && result.Title.RequiresPlatformClient() && !result.PlatformClientInstalled {
				message = fmt.Sprintf("installed, but required platform client is missing (%s)", result.Title.PlatformClient.Platform)
			} else if result.PreviouslyRegistered {
				message = "launcher already registered"
			} else {
				message = "launcher registered successfully"
			}
			fmt.Printf("%s: %s\n", result.Title.GameLabel, message)
		}
	} else if len(args) == 1 {
		err := gameRouter.StartGame(args[0])
		if err != nil {
			fmt.Printf("Failed to launch based on URL: %s (%s)\n", args[0], err)
		} else {
			fmt.Printf("Launched game based on URL: %s\n", args[0])
		}
	}
	fmt.Println("Window will close in 15 seconds")
	time.Sleep(15 * time.Second)
}
