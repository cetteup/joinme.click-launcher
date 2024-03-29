package router

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/internal/domain"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

type action int
type actionPathKey string

const (
	regPathSoftware         = "SOFTWARE"
	regPathClasses          = "Classes"
	regPathOpen             = "open"
	regPathShell            = "shell"
	regPathCommand          = "command"
	regValueNameDefault     = ""
	regValueNameURLProtocol = "URL Protocol"

	actionUrlHostname = "act"

	actionLaunchAndJoin action = iota
	actionLaunchOnly

	actionPathKeyLaunch actionPathKey = "launch"
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
	DeleteKey(k registry.Key, path string) error
}

type GameFinder interface {
	IsInstalledAnywhere(configs []software_finder.Config) (bool, error)
	IsInstalled(config software_finder.Config) (bool, error)
	GetInstallDirFromSomewhere(configs []software_finder.Config) (string, error)
	GetInstallDir(config software_finder.Config) (string, error)
}

type GameLauncher interface {
	StartGame(u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType, cmdBuilder game_launcher.CommandBuilder, hookHandlers ...game_launcher.HookHandler) error
}

type GameRouter struct {
	repository RegistryRepository
	finder     GameFinder
	launcher   GameLauncher
	GameTitles map[string]domain.GameTitle
}

type handlerRegistrationResult struct {
	Title                   domain.GameTitle
	GameInstalled           bool
	PlatformClientInstalled bool
	PreviouslyRegistered    bool
	Registered              bool
	Error                   error
}

func New(repository RegistryRepository, finder GameFinder, launcher GameLauncher) *GameRouter {
	return &GameRouter{
		repository: repository,
		finder:     finder,
		launcher:   launcher,
		GameTitles: map[string]domain.GameTitle{},
	}
}

func (r *GameRouter) AddTitle(gameTitles ...domain.GameTitle) {
	for _, gt := range gameTitles {
		customConfig := internal.Config.GetCustomLauncherConfig(gt.ProtocolScheme)
		if customConfig.HasValues() {
			gt.AddCustomConfig(*customConfig)
		}
		r.GameTitles[gt.ProtocolScheme] = gt
	}
}

func (r *GameRouter) RegisterHandlers() []handlerRegistrationResult {
	results := make([]handlerRegistrationResult, 0, len(r.GameTitles))
	for _, gameTitle := range r.GameTitles {
		result := handlerRegistrationResult{
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

		registered, err := r.isHandlerRegistered(gameTitle)
		if err != nil {
			result.Error = fmt.Errorf("failed to determine whether handler is registered: %e", err)
			results = append(results, result)
			continue
		}
		result.PreviouslyRegistered = registered

		if !registered {
			if err = r.registerHandler(gameTitle); err != nil {
				result.Error = fmt.Errorf("failed to register as URL protocol handler: %e", err)
			} else {
				result.Registered = true
			}
		}

		results = append(results, result)
	}

	return results
}

func (r *GameRouter) DeregisterHandlers() error {
	for _, gameTitle := range r.GameTitles {
		err := r.deregisterHandler(gameTitle)
		if err != nil {
			return fmt.Errorf("failed to deregister handler for %s (%s): %w", gameTitle.Name, gameTitle.ProtocolScheme, err)
		}
	}
	return nil
}

func (r *GameRouter) isHandlerRegistered(gameTitle domain.GameTitle) (bool, error) {
	path := r.getUrlHandlerRegistryPath(gameTitle, regPathShell, regPathOpen, regPathCommand)
	value, err := r.repository.GetStringValue(registry.CURRENT_USER, path, regValueNameDefault)
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

func (r *GameRouter) registerHandler(gameTitle domain.GameTitle) error {
	basePath := r.getUrlHandlerRegistryPath(gameTitle)
	err := r.repository.CreateKey(registry.CURRENT_USER, basePath)
	if err != nil {
		return err
	}

	err = r.repository.SetStringValue(registry.CURRENT_USER, basePath, regValueNameDefault, fmt.Sprintf("URL:%s protocol", gameTitle.Name))
	if err != nil {
		return err
	}
	err = r.repository.SetStringValue(registry.CURRENT_USER, basePath, regValueNameURLProtocol, "")
	if err != nil {
		return err
	}

	subKeys := []string{regPathShell, regPathOpen, regPathCommand}
	for i := range subKeys {
		subPath := r.getUrlHandlerRegistryPath(gameTitle, subKeys[:i+1]...)
		err = r.repository.CreateKey(registry.CURRENT_USER, subPath)
		if err != nil {
			return err
		}
	}

	cmdPath := r.getUrlHandlerRegistryPath(gameTitle, subKeys...)
	cmd, err := r.getHandlerCommand()
	if err != nil {
		return err
	}

	return r.repository.SetStringValue(registry.CURRENT_USER, cmdPath, regValueNameDefault, cmd)
}

func (r *GameRouter) deregisterHandler(gameTitle domain.GameTitle) error {
	keys := []string{r.getUrlHandlerRegistryPath(gameTitle)}
	subKeys := []string{regPathShell, regPathOpen, regPathCommand}
	for i := range subKeys {
		// We need to delete keys in reverse/"descending" order, so prepend to list
		keys = append([]string{r.getUrlHandlerRegistryPath(gameTitle, subKeys[:i+1]...)}, keys...)
	}

	for _, key := range keys {
		err := r.repository.DeleteKey(registry.CURRENT_USER, key)
		if err != nil && !errors.Is(err, registry.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (r *GameRouter) getUrlHandlerRegistryPath(gameTitle domain.GameTitle, children ...string) string {
	path := filepath.Join(regPathSoftware, regPathClasses, gameTitle.ProtocolScheme)
	for _, child := range children {
		path = filepath.Join(path, child)
	}

	return path
}

func (r *GameRouter) getHandlerCommand() (string, error) {
	launcherPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" \"%%1\"", launcherPath), nil
}

func (r *GameRouter) RunURL(commandLineUrl string) (*domain.GameTitle, error) {
	u, err := url.Parse(commandLineUrl)
	if err != nil {
		return nil, err
	}

	gameTitle, ok := r.GameTitles[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("game not supported: %s", u.Scheme)
	}

	if err = r.ensurePrerequisites(gameTitle, u); err != nil {
		return &gameTitle, err
	}

	action, err := r.getActionFromURL(u)
	if err != nil {
		return &gameTitle, err
	}

	switch action {
	case actionLaunchOnly:
		return &gameTitle, r.startGame(gameTitle, u, game_launcher.LaunchTypeLaunchOnly)
	default:
		return &gameTitle, r.startGame(gameTitle, u, game_launcher.LaunchTypeLaunchAndJoin)
	}
}

func (r *GameRouter) ensurePrerequisites(gameTitle domain.GameTitle, u *url.URL) error {
	if err := r.ensureGameIsInstalled(gameTitle); err != nil {
		return err
	}
	if err := r.ensurePlatformClientIsInstalledIfRequired(gameTitle); err != nil {
		return err
	}
	if err := r.ensureModIsSupportedAndInstalledIfGiven(gameTitle, u); err != nil {
		return err
	}

	return nil
}

func (r *GameRouter) ensureGameIsInstalled(gameTitle domain.GameTitle) error {
	gameInstalled, err := r.finder.IsInstalledAnywhere(gameTitle.FinderConfigs)
	if err != nil {
		return err
	}
	if !gameInstalled {
		return fmt.Errorf("game not installed")
	}

	return nil
}

func (r *GameRouter) ensurePlatformClientIsInstalledIfRequired(gameTitle domain.GameTitle) error {
	if !gameTitle.RequiresPlatformClient() {
		return nil
	}

	clientInstalled, err := r.finder.IsInstalled(gameTitle.PlatformClient.FinderConfig)
	if err != nil {
		return err
	}
	if !clientInstalled {
		return fmt.Errorf("required platform client not installed: %s", gameTitle.PlatformClient.Platform)
	}

	return nil
}

func (r *GameRouter) ensureModIsSupportedAndInstalledIfGiven(gameTitle domain.GameTitle, u *url.URL) error {
	query := u.Query()
	if !internal.QueryHasMod(query) {
		return nil
	}

	slug := internal.GetModFromQuery(query)
	mod := gameTitle.GetMod(slug)
	if mod == nil {
		return fmt.Errorf("mod not supported: %s", slug)
	}

	gameInstallPath, err := r.finder.GetInstallDirFromSomewhere(gameTitle.FinderConfigs)
	if err != nil {
		return err
	}
	modInstalled, err := r.finder.IsInstalledAnywhere(mod.ComputeFinderConfigs(gameInstallPath))
	if err != nil {
		return err
	}
	if !modInstalled {
		return fmt.Errorf("mod not installed: %s", mod.Name)
	}

	return nil
}

func (r *GameRouter) startGame(gameTitle domain.GameTitle, u *url.URL, launchType game_launcher.LaunchType) error {
	// Only join url use/require URL parameters, so only validate those
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		err := gameTitle.URLValidator.Validate(u)
		if err != nil {
			return err
		}
	}

	// Build final launcher config
	launcherConfig, err := r.buildLauncherConfig(gameTitle)
	if err != nil {
		return err
	}

	return r.launcher.StartGame(u, launcherConfig, launchType, gameTitle.CmdBuilder, gameTitle.HookHandlers...)
}

func (r *GameRouter) buildLauncherConfig(gameTitle domain.GameTitle) (game_launcher.Config, error) {
	launcherConfig := gameTitle.LauncherConfig
	if gameTitle.RequiresPlatformClient() {
		launcherConfig.ExecutableName = gameTitle.PlatformClient.LauncherConfig.ExecutableName
		launcherConfig.ExecutablePath = gameTitle.PlatformClient.LauncherConfig.ExecutablePath
		installPath, err := r.finder.GetInstallDir(gameTitle.PlatformClient.FinderConfig)
		if err != nil {
			return game_launcher.Config{}, err
		}
		launcherConfig.InstallPath = installPath
	} else {
		installPath, err := r.finder.GetInstallDirFromSomewhere(gameTitle.FinderConfigs)
		if err != nil {
			return game_launcher.Config{}, err
		}
		launcherConfig.InstallPath = installPath
	}

	return launcherConfig, nil
}

func (r *GameRouter) getActionFromURL(u *url.URL) (action, error) {
	if !r.isActionURL(u) {
		return actionLaunchAndJoin, nil
	}

	switch strings.TrimPrefix(u.Path, "/") {
	case string(actionPathKeyLaunch):
		return actionLaunchOnly, nil
	default:
		return 0, fmt.Errorf("action not supported: %s", u.Path)
	}
}

func (r *GameRouter) isActionURL(u *url.URL) bool {
	return u.Hostname() == actionUrlHostname
}
