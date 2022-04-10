//go:generate goversioninfo

package main

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
	"github.com/cetteup/joinme.click-launcher/game/titles"
	"github.com/cetteup/joinme.click-launcher/internal"
	"os"
	"sort"
	"time"
)

func init() {
	registryRepository := internal.NewRegistryRepository()
	err := internal.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load configuration: %s\n", err)
		fmt.Println("Continuing with defaults")
	}

	router = game.NewRouter(registryRepository)
	router.AddLauncher(titles.Bf1942Config, titles.Bf1942CmdBuilder)
	router.AddLauncher(titles.BfVietnamConfig, titles.BfVietnamCmdBuilder)
	router.AddLauncher(titles.Bf2Config, titles.Bf2CmdBuilder)
	router.AddLauncher(titles.Bf2SFConfig, titles.Bf2CmdBuilder)
	router.AddLauncher(titles.CodConfig, game.PlusConnectCmdBuilder)
	router.AddLauncher(titles.CodUOConfig, game.PlusConnectCmdBuilder)
	router.AddLauncher(titles.Cod2Config, game.PlusConnectCmdBuilder)
	router.AddLauncher(titles.Cod4Config, game.PlusConnectCmdBuilder)
	router.AddLauncher(titles.CodWawConfig, game.PlusConnectCmdBuilder)
	router.AddLauncher(titles.FearSec2Config, titles.FearSec2CmdBuilder)
	router.AddLauncher(titles.ParaworldConfig, titles.ParaworldCmdBuilder)
	router.AddLauncher(titles.Swat4Config, titles.Swat4CmdBuilder)
	router.AddLauncher(titles.Swat4XConfig, titles.Swat4CmdBuilder)
	router.AddLauncher(titles.VietcongConfig, titles.VietcongCmdBuilder)
}

var (
	router *game.Router
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		results := router.RegisterHandlers()
		sort.Slice(results, func(i, j int) bool {
			return results[i].ProtocolScheme < results[j].ProtocolScheme
		})
		for _, result := range results {
			var message string
			if result.Error != nil {
				message = fmt.Sprintf("handler registration failed (%s)\n", result.Error)
			} else if !result.Installed {
				message = "not installed"
			} else if result.PreviouslyRegistered {
				message = "launcher already registered"
			} else {
				message = "launcher registered successfully"
			}
			fmt.Printf("%s: %s\n", result.GameLabel, message)
		}
	} else if len(args) == 1 {
		err := router.StartGame(args[0])
		if err != nil {
			fmt.Printf("Failed to launch based on URL: %s (%s)\n", args[0], err)
		} else {
			fmt.Printf("Launched game based on URL: %s\n", args[0])
		}
	}
	fmt.Println("Window will close in 15 seconds")
	time.Sleep(15 * time.Second)
}
