//go:generate goversioninfo

package main

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
	"github.com/cetteup/joinme.click-launcher/game/titles"
	"github.com/cetteup/joinme.click-launcher/internal"
	"os"
	"time"
)

func init() {
	repository := internal.NewRegistryRepository()
	router = game.NewRouter(repository)
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
		for _, gameLauncher := range router.Launchers {
			installed, err := gameLauncher.IsGameInstalled()
			if err != nil {
				fmt.Println(fmt.Errorf("failed to determine whether %s is installed: %e", gameLauncher.Config.GameLabel, err))
				continue
			}

			if !installed {
				fmt.Printf("%s is not installed according to the registry\n", gameLauncher.Config.GameLabel)
				continue
			}

			registered, err := gameLauncher.IsHandlerRegistered()
			if err != nil {
				fmt.Println(fmt.Errorf("failed to determine whether %s handler is registered: %e", gameLauncher.Config.GameLabel, err))
				continue
			}

			if !registered {
				fmt.Printf("Detected %s install, registering launcher for %s protocol\n", gameLauncher.Config.GameLabel, gameLauncher.Config.ProtocolScheme)
				if err := gameLauncher.RegisterHandler(); err != nil {
					fmt.Println(fmt.Errorf("failed to register as URL protocol handler for %s: %e", gameLauncher.Config.GameLabel, err))
				}
			} else {
				fmt.Printf("Detected %s install, launcher already registered for %s protocol\n", gameLauncher.Config.GameLabel, gameLauncher.Config.ProtocolScheme)
			}
		}
	} else if len(args) == 1 {
		err := router.StartGame(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		// TODO Remove later
		fmt.Printf("Launched game based on URL: %s\n", args[0])
	}
	fmt.Println("Window will close in 15 seconds")
	time.Sleep(15 * time.Second)
}
