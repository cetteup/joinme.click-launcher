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
	cod := game.NewLauncher(repository, titles.CodConfig, game.PlusConnectCmdBuilder)
	codUO := game.NewLauncher(repository, titles.CodUOConfig, game.PlusConnectCmdBuilder)
	cod2 := game.NewLauncher(repository, titles.Cod2Config, game.PlusConnectCmdBuilder)
	cod4 := game.NewLauncher(repository, titles.Cod4Config, game.PlusConnectCmdBuilder)
	codWaw := game.NewLauncher(repository, titles.CodWawConfig, game.PlusConnectCmdBuilder)
	fearSec2 := game.NewLauncher(repository, titles.FearSec2Config, titles.FearSec2CmdBuilder)
	paraworld := game.NewLauncher(repository, titles.ParaworldConfig, titles.ParaworldCmdBuilder)
	swat4 := game.NewLauncher(repository, titles.Swat4Config, titles.Swat4CmdBuilder)
	swat4x := game.NewLauncher(repository, titles.Swat4XConfig, titles.Swat4CmdBuilder)
	vietcong := game.NewLauncher(repository, titles.VietcongConfig, titles.VietcongCmdBuilder)

	launchers = map[string]*game.Launcher{
		bf1942.Config.ProtocolScheme:    bf1942,
		bfVietnam.Config.ProtocolScheme: bfVietnam,
		bf2.Config.ProtocolScheme:       bf2,
		cod.Config.ProtocolScheme:       cod,
		codUO.Config.ProtocolScheme:     codUO,
		cod2.Config.ProtocolScheme:      cod2,
		cod4.Config.ProtocolScheme:      cod4,
		codWaw.Config.ProtocolScheme:    codWaw,
		fearSec2.Config.ProtocolScheme:  fearSec2,
		paraworld.Config.ProtocolScheme: paraworld,
		swat4.Config.ProtocolScheme:     swat4,
		swat4x.Config.ProtocolScheme:    swat4x,
		vietcong.Config.ProtocolScheme:  vietcong,
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
