//go:generate goversioninfo

package main

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game"
	"github.com/cetteup/joinme.click-launcher/game/titles"
	"github.com/cetteup/joinme.click-launcher/internal"
	"net/url"
	"os"
	"time"
)

func init() {
	repository := internal.NewRegistryRepository()
	bf1942 := game.NewLauncher(repository, titles.Bf1942Config, titles.Bf1942CmdBuilder)
	bfVietnam := game.NewLauncher(repository, titles.BfVietnamConfig, titles.BfVietnamCmdBuilder)
	bf2 := game.NewLauncher(repository, titles.Bf2Config, titles.Bf2CmdBuilder)
	codWaw := game.NewLauncher(repository, titles.CodWawConfig, titles.CodWawCmdBuilder)
	fearSec2 := game.NewLauncher(repository, titles.FearSec2Config, titles.FearSec2CmdBuilder)
	paraworld := game.NewLauncher(repository, titles.ParaworldConfig, titles.ParaworldCmdBuilder)

	launchers = map[string]*game.Launcher{
		"bf1942":    bf1942,
		"bfvietnam": bfVietnam,
		"bf2":       bf2,
		"codwaw":    codWaw,
		"fearsec2":  fearSec2,
		"paraworld": paraworld,
	}
}

var (
	launchers map[string]*game.Launcher
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		for _, gameLauncher := range launchers {
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

			if installed && !registered {
				fmt.Printf("Detected %s install, registering launcher for %s protocol\n", gameLauncher.Config.GameLabel, gameLauncher.Config.ProtocolScheme)
				if err := gameLauncher.RegisterHandler(); err != nil {
					fmt.Println(fmt.Errorf("failed to register as URL protocol handler for %s: %e", gameLauncher.Config.GameLabel, err))
				}
			} else if installed {
				fmt.Printf("Detected %s install, launcher already registered for %s protocol\n", gameLauncher.Config.GameLabel, gameLauncher.Config.ProtocolScheme)
			}
		}
	} else if len(args) == 1 {
		u, err := url.Parse(args[0])
		if err != nil {
			panic(err)
		}

		gameLauncher, ok := launchers[u.Scheme]
		if !ok {
			panic("Game not supported")
		}

		port := u.Port()
		if port == "" {
			panic("No port given in URL")
		}

		err = gameLauncher.StartGame(u.Hostname(), port)
		if err != nil {
			panic(err)
		}

		// TODO Remove later
		fmt.Printf("Launching %s to join %s:%s\n", gameLauncher.Config.GameLabel, u.Hostname(), port)
	}
	fmt.Println("Window will close in 15 seconds")
	time.Sleep(15 * time.Second)
}
