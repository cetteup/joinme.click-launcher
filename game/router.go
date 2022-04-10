package game

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/internal"
	"net/url"
)

type Router struct {
	repository RegistryRepository
	Launchers  map[string]*Launcher
}

type HandlerRegistrationResult struct {
	ProtocolScheme       string
	GameLabel            string
	Installed            bool
	PreviouslyRegistered bool
	Registered           bool
	Error                error
}

func NewRouter(repository RegistryRepository) *Router {
	return &Router{
		repository: repository,
		Launchers:  map[string]*Launcher{},
	}
}

func (r Router) AddLauncher(config LauncherConfig, cmdBuilder CommandBuilder) {
	config.Custom = internal.RunningConfig.GetCustomLauncherConfig(config.ProtocolScheme)
	launcher := NewLauncher(r.repository, config, cmdBuilder)
	r.Launchers[launcher.Config.ProtocolScheme] = launcher
}

func (r Router) RegisterHandlers() []HandlerRegistrationResult {
	results := make([]HandlerRegistrationResult, 0, len(r.Launchers))
	for _, launcher := range r.Launchers {
		result := HandlerRegistrationResult{
			ProtocolScheme: launcher.Config.ProtocolScheme,
			GameLabel:      launcher.Config.GameLabel,
		}

		installed, err := launcher.IsGameInstalled()
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

		registered, err := launcher.IsHandlerRegistered()
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether handler is registered: %e", err)
			results = append(results, result)
			continue
		}
		result.PreviouslyRegistered = registered

		if !registered {
			if err = launcher.RegisterHandler(); err != nil {
				result.Error = fmt.Errorf("failed to register as URL protocol handler: %e", err)
			} else {
				result.Registered = true
			}
		}

		results = append(results, result)
	}

	return results
}

func (r Router) StartGame(commandLineUrl string) error {
	u, err := url.Parse(commandLineUrl)
	if err != nil {
		return err
	}

	launcher, ok := r.Launchers[u.Scheme]
	if !ok {
		return fmt.Errorf("game not supported: %s", u.Scheme)
	}

	port := u.Port()
	if port == "" {
		return fmt.Errorf("no port given in URL: %s", port)
	}

	err = launcher.StartGame(u.Hostname(), port)
	if err != nil {
		return err
	}

	return nil
}
