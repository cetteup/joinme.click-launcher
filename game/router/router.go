package router

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"github.com/cetteup/joinme.click-launcher/game/title"
	"github.com/cetteup/joinme.click-launcher/internal"
	"golang.org/x/sys/windows/registry"
)

const (
	RegPathSoftware   = "SOFTWARE"
	RegPathClasses    = "Classes"
	RegPathOpen       = "open"
	RegPathShell      = "shell"
	RegPathCommand    = "command"
	RegKeyDefault     = ""
	RegKeyURLProtocol = "URL Protocol"
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type GameFinder interface {
	IsInstalledAnywhere(configs []finder.Config) (bool, error)
	IsInstalled(config finder.Config) (bool, error)
	GetInstallDirFromSomewhere(configs []finder.Config) (string, error)
	GetInstallDir(config finder.Config) (string, error)
}

type GameRouter struct {
	repository RegistryRepository
	finder     GameFinder
	GameTitles map[string]title.GameTitle
}

type HandlerRegistrationResult struct {
	Title                   title.GameTitle
	GameInstalled           bool
	PlatformClientInstalled bool
	PreviouslyRegistered    bool
	Registered              bool
	Error                   error
}

func NewGameRouter(repository RegistryRepository, finder GameFinder) *GameRouter {
	return &GameRouter{
		repository: repository,
		finder:     finder,
		GameTitles: map[string]title.GameTitle{},
	}
}

func (r GameRouter) AddTitle(gameTitle title.GameTitle) {
	customConfig := internal.Config.GetCustomLauncherConfig(gameTitle.ProtocolScheme)
	if customConfig.HasValues() {
		gameTitle.AddCustomConfig(*customConfig)
	}
	r.GameTitles[gameTitle.ProtocolScheme] = gameTitle
}

func (r GameRouter) RegisterHandlers() []HandlerRegistrationResult {
	results := make([]HandlerRegistrationResult, 0, len(r.GameTitles))
	for _, gameTitle := range r.GameTitles {
		result := HandlerRegistrationResult{
			Title: gameTitle,
		}

		installed, err := r.finder.IsInstalledAnywhere(gameTitle.FinderConfigs)
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether game is installed: %e", err)
			results = append(results, result)
			continue
		}
		result.GameInstalled = installed

		if !installed {
			results = append(results, result)
			continue
		}

		if gameTitle.RequiresPlatformClient() {
			platformClientInstalled, err := r.finder.IsInstalled(gameTitle.PlatformClient.FinderConfig)
			if err != nil {
				result.Error = fmt.Errorf("failed to determine whether required platform (%s) is installed: %e", gameTitle.PlatformClient.Platform, err)
				results = append(results, result)
				continue
			}
			result.PlatformClientInstalled = platformClientInstalled

			if !platformClientInstalled {
				results = append(results, result)
				continue
			}
		}

		registered, err := r.IsHandlerRegistered(gameTitle)
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether handler is registered: %e", err)
			results = append(results, result)
			continue
		}
		result.PreviouslyRegistered = registered

		if !registered {
			if err = r.RegisterHandler(gameTitle); err != nil {
				result.Error = fmt.Errorf("failed to register as URL protocol handler: %e", err)
			} else {
				result.Registered = true
			}
		}

		results = append(results, result)
	}

	return results
}

func (r GameRouter) IsHandlerRegistered(gameTitle title.GameTitle) (bool, error) {
	path := r.getUrlHandlerRegistryPath(gameTitle, []string{RegPathShell, RegPathOpen, RegPathCommand})
	value, err := r.repository.GetStringValue(registry.CURRENT_USER, path, RegKeyDefault)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	expected, err := r.getHandlerCommand()
	if err != nil {
		return false, err
	}

	return value == expected, nil
}

func (r GameRouter) RegisterHandler(gameTitle title.GameTitle) error {
	basePath := r.getUrlHandlerRegistryPath(gameTitle, nil)
	err := r.repository.CreateKey(registry.CURRENT_USER, basePath)
	if err != nil {
		return err
	}

	err = r.repository.SetStringValue(registry.CURRENT_USER, basePath, RegKeyDefault, fmt.Sprintf("URL:%s protocol", gameTitle.GameLabel))
	if err != nil {
		return err
	}
	err = r.repository.SetStringValue(registry.CURRENT_USER, basePath, RegKeyURLProtocol, "")
	if err != nil {
		return err
	}

	subKeys := []string{RegPathShell, RegPathOpen, RegPathCommand}
	for i := range subKeys {
		subPath := r.getUrlHandlerRegistryPath(gameTitle, subKeys[:i+1])
		err = r.repository.CreateKey(registry.CURRENT_USER, subPath)
		if err != nil {
			return err
		}
	}

	cmdPath := r.getUrlHandlerRegistryPath(gameTitle, subKeys)
	cmd, err := r.getHandlerCommand()
	if err != nil {
		return err
	}

	return r.repository.SetStringValue(registry.CURRENT_USER, cmdPath, RegKeyDefault, cmd)
}

func (r GameRouter) getUrlHandlerRegistryPath(gameTitle title.GameTitle, children []string) string {
	path := filepath.Join(RegPathSoftware, RegPathClasses, gameTitle.ProtocolScheme)
	for _, child := range children {
		path = filepath.Join(path, child)
	}

	return path
}

func (r GameRouter) getHandlerCommand() (string, error) {
	launcherPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" \"%%1\"", launcherPath), nil
}

func (r GameRouter) StartGame(commandLineUrl string) error {
	u, err := url.Parse(commandLineUrl)
	if err != nil {
		return err
	}

	gameTitle, ok := r.GameTitles[u.Scheme]
	if !ok {
		return fmt.Errorf("game not supported: %s", u.Scheme)
	}

	port := u.Port()
	if gameTitle.RequiresPort && port == "" {
		return fmt.Errorf("port is required but was not given in URL")
	}

	// Build final launcher config
	launcherConfig := gameTitle.LauncherConfig
	if gameTitle.RequiresPlatformClient() {
		launcherConfig.ExecutablePath = gameTitle.PlatformClient.LauncherConfig.ExecutablePath
		launcherConfig.InstallPath, err = r.finder.GetInstallDir(gameTitle.PlatformClient.FinderConfig)
	} else {
		launcherConfig.InstallPath, err = r.finder.GetInstallDirFromSomewhere(gameTitle.FinderConfigs)
	}
	if err != nil {
		return err
	}

	gameLauncher := launcher.NewGameLauncher(launcherConfig, gameTitle.CmdBuilder)
	err = gameLauncher.StartGame(u.Scheme, u.Hostname(), port, u)
	if err != nil {
		return err
	}

	return nil
}
