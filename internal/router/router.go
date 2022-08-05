package router

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
	"golang.org/x/sys/windows/registry"
)

type Action int
type ActionPathKey string

const (
	RegPathSoftware   = "SOFTWARE"
	RegPathClasses    = "Classes"
	RegPathOpen       = "open"
	RegPathShell      = "shell"
	RegPathCommand    = "command"
	RegKeyDefault     = ""
	RegKeyURLProtocol = "URL Protocol"

	ActionUrlHostname = "act"

	ActionLaunchAndJoin Action = iota
	ActionLaunchOnly

	ActionPathKeyLaunch ActionPathKey = "launch"
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type GameFinder interface {
	IsInstalledAnywhere(configs []software_finder.Config) (bool, error)
	IsInstalled(config software_finder.Config) (bool, error)
	GetInstallDirFromSomewhere(configs []software_finder.Config) (string, error)
	GetInstallDir(config software_finder.Config) (string, error)
}

type GameRouter struct {
	repository RegistryRepository
	finder     GameFinder
	GameTitles map[string]domain.GameTitle
}

type HandlerRegistrationResult struct {
	Title                   domain.GameTitle
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
		GameTitles: map[string]domain.GameTitle{},
	}
}

func (r GameRouter) AddTitle(gameTitle domain.GameTitle) {
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

func (r GameRouter) IsHandlerRegistered(gameTitle domain.GameTitle) (bool, error) {
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

func (r GameRouter) RegisterHandler(gameTitle domain.GameTitle) error {
	basePath := r.getUrlHandlerRegistryPath(gameTitle, nil)
	err := r.repository.CreateKey(registry.CURRENT_USER, basePath)
	if err != nil {
		return err
	}

	err = r.repository.SetStringValue(registry.CURRENT_USER, basePath, RegKeyDefault, fmt.Sprintf("URL:%s protocol", gameTitle.Name))
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

func (r GameRouter) getUrlHandlerRegistryPath(gameTitle domain.GameTitle, children []string) string {
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

func (r GameRouter) RunURL(commandLineUrl string) error {
	u, err := url.Parse(commandLineUrl)
	if err != nil {
		return err
	}

	gameTitle, ok := r.GameTitles[u.Scheme]
	if !ok {
		return fmt.Errorf("game not supported: %s", u.Scheme)
	}

	action, err := r.getActionFromURL(u)
	if err != nil {
		return err
	}

	switch action {
	case ActionLaunchOnly:
		return r.startGame(gameTitle, u, game_launcher.LaunchTypeLaunchOnly)
	default:
		return r.startGame(gameTitle, u, game_launcher.LaunchTypeLaunchAndJoin)
	}
}

func (r GameRouter) startGame(gameTitle domain.GameTitle, u *url.URL, launchType game_launcher.LaunchType) error {
	// Only join url use/require URL parameters, so only validate those
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		err := gameTitle.URLValidator(u)
		if err != nil {
			return err
		}
	}

	// Build final launcher config
	launcherConfig := gameTitle.LauncherConfig
	if gameTitle.RequiresPlatformClient() {
		launcherConfig.ExecutableName = gameTitle.PlatformClient.LauncherConfig.ExecutableName
		launcherConfig.ExecutablePath = gameTitle.PlatformClient.LauncherConfig.ExecutablePath
		installPath, err := r.finder.GetInstallDir(gameTitle.PlatformClient.FinderConfig)
		if err != nil {
			return err
		}
		launcherConfig.InstallPath = installPath
	} else {
		installPath, err := r.finder.GetInstallDirFromSomewhere(gameTitle.FinderConfigs)
		if err != nil {
			return err
		}
		launcherConfig.InstallPath = installPath
	}

	// Always use the game launcher.Config for preparation, since we need to (for example) kill the game, not the platform client before launch
	if err := game_launcher.PrepareLaunch(gameTitle.LauncherConfig); err != nil {
		return err
	}
	if err := game_launcher.StartGame(u, launcherConfig, launchType, gameTitle.CmdBuilder); err != nil {
		return err
	}

	return nil
}

func (r GameRouter) getActionFromURL(u *url.URL) (Action, error) {
	if !r.isActionURL(u) {
		return ActionLaunchAndJoin, nil
	}

	switch strings.TrimPrefix(u.Path, "/") {
	case string(ActionPathKeyLaunch):
		return ActionLaunchOnly, nil
	default:
		return 0, fmt.Errorf("action not supported: %s", u.Path)
	}
}

func (r GameRouter) isActionURL(u *url.URL) bool {
	return u.Hostname() == ActionUrlHostname
}
