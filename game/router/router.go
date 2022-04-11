package router

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/internal"
	"golang.org/x/sys/windows/registry"
	"net/url"
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type GameRouter struct {
	repository    RegistryRepository
	GameLaunchers map[string]*launcher.GameLauncher
}

type HandlerRegistrationResult struct {
	ProtocolScheme       string
	GameLabel            string
	Installed            bool
	PreviouslyRegistered bool
	Registered           bool
	Error                error
}

func NewGameRouter(repository RegistryRepository) *GameRouter {
	return &GameRouter{
		repository:    repository,
		GameLaunchers: map[string]*launcher.GameLauncher{},
	}
}

func (r GameRouter) AddLauncher(config launcher.Config, cmdBuilder launcher.CommandBuilder) {
	config.Custom = internal.RunningConfig.GetCustomLauncherConfig(config.ProtocolScheme)
	gameLauncher := launcher.NewGameLauncher(r.repository, config, cmdBuilder)
	r.GameLaunchers[gameLauncher.Config.ProtocolScheme] = gameLauncher
}

func (r GameRouter) RegisterHandlers() []HandlerRegistrationResult {
	results := make([]HandlerRegistrationResult, 0, len(r.GameLaunchers))
	for _, gameLauncher := range r.GameLaunchers {
		result := HandlerRegistrationResult{
			ProtocolScheme: gameLauncher.Config.ProtocolScheme,
			GameLabel:      gameLauncher.Config.GameLabel,
		}

		installed, err := gameLauncher.IsGameInstalled()
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether game is installed: %e", err)
			results = append(results, result)
			continue
		}
		result.Installed = installed

		if !installed {
			results = append(results, result)
			continue
		}

		registered, err := gameLauncher.IsHandlerRegistered()
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether handler is registered: %e", err)
			results = append(results, result)
			continue
		}
		result.PreviouslyRegistered = registered

		if !registered {
			if err = gameLauncher.RegisterHandler(); err != nil {
				result.Error = fmt.Errorf("failed to register as URL protocol handler: %e", err)
			} else {
				result.Registered = true
			}
		}

		results = append(results, result)
	}

	return results
}

func (r GameRouter) StartGame(commandLineUrl string) error {
	u, err := url.Parse(commandLineUrl)
	if err != nil {
		return err
	}

	gameLauncher, ok := r.GameLaunchers[u.Scheme]
	if !ok {
		return fmt.Errorf("game not supported: %s", u.Scheme)
	}

	port := u.Port()
	if port == "" {
		return fmt.Errorf("no port given in URL: %s", port)
	}

	err = gameLauncher.StartGame(u.Hostname(), port)
	if err != nil {
		return err
	}

	return nil
}
