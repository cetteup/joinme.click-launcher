package game

import (
	"fmt"
	"net/url"
)

type Router struct {
	repository RegistryRepository
	Launchers  map[string]*Launcher
}

func NewRouter(repository RegistryRepository) *Router {
	return &Router{
		repository: repository,
		Launchers:  map[string]*Launcher{},
	}
}

func (r Router) AddLauncher(config LauncherConfig, cmdBuilder CommandBuilder) {
	launcher := NewLauncher(r.repository, config, cmdBuilder)
	r.Launchers[launcher.Config.ProtocolScheme] = launcher
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
