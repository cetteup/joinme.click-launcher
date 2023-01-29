package game_launcher

import (
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type LaunchDir int
type LaunchType int
type HookWhen string

const (
	LaunchDirInstallDir LaunchDir = iota
	LaunchDirBinaryDir

	LaunchTypeLaunchAndJoin LaunchType = iota
	LaunchTypeLaunchOnly

	HookWhenAlways     HookWhen = "always"
	HookWhenPreLaunch  HookWhen = "pre-launch"
	HookWhenPostLaunch HookWhen = "post-launch"

	handlerLogKey = "handler"
)

type FileRepository interface {
	FileExists(path string) (bool, error)
	DirExists(path string) (bool, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]os.DirEntry, error)
	Glob(pattern string) ([]string, error)
	RemoveAll(path string) error
}

type GameLauncher struct {
	repository FileRepository
}

func New(repository FileRepository) *GameLauncher {
	return &GameLauncher{
		repository: repository,
	}
}

type Config struct {
	DefaultArgs       []string
	AppendDefaultArgs bool
	StartIn           LaunchDir
	ExecutableName    string
	// Relative path from install path to folder containing the executable
	ExecutablePath string
	InstallPath    string
	HookConfigs    []HookConfig
}

type HookConfig struct {
	Handler     string
	When        HookWhen
	ExitOnError bool
	Args        map[string]string
}

// TODO Interface should either move to router or game launcher should do the URL validation
type URLValidator interface {
	Validate(u *url.URL) error
}

type CommandBuilder interface {
	// GetArgs Construct slice of launch arguments for a game. Receives a file repository to be able to access any
	// config file it may need to construct the arguments.
	GetArgs(fr FileRepository, u *url.URL, config Config, launchType LaunchType) ([]string, error)
}

type HookHandler interface {
	Run(fr FileRepository, u *url.URL, config Config, launchType LaunchType, args map[string]string) error
	String() string
}

func (l *GameLauncher) StartGame(u *url.URL, config Config, launchType LaunchType, cmdBuilder CommandBuilder, hookHandlers ...HookHandler) error {
	// Convert handlers to map to make access faster/easier
	hookHandlerMap := toHookHandlerMap(hookHandlers)

	// Run pre-launch hooks
	if err := l.runHooks(u, config, launchType, hookHandlerMap, HookWhenPreLaunch); err != nil {
		return err
	}

	// Start the game
	if err := l.startGame(u, config, launchType, cmdBuilder); err != nil {
		return err
	}

	// Run post-launch hooks
	return l.runHooks(u, config, launchType, hookHandlerMap, HookWhenPostLaunch)
}

func (l *GameLauncher) runHooks(u *url.URL, config Config, launchType LaunchType, handlers map[string]HookHandler, when HookWhen) error {
	for _, hc := range config.HookConfigs {
		if hc.When != when && hc.When != HookWhenAlways {
			log.Debug().Str(handlerLogKey, hc.Handler).Str("when", string(when)).Msg("Skipping hook handler not configured to run now")
			continue
		}

		handler, ok := handlers[hc.Handler]
		if !ok {
			log.Warn().Str(handlerLogKey, hc.Handler).Msg("Skipping unknown hook handler")
			continue
		}

		log.Debug().Str(handlerLogKey, hc.Handler).Interface("args", hc.Args).Msg("Running hook handler")

		err := handler.Run(l.repository, u, config, launchType, hc.Args)
		if err != nil {
			log.Error().Err(err).Str(handlerLogKey, hc.Handler).Msg("Hook handler execution failed")
			if hc.ExitOnError {
				return err
			}
		}
	}
	return nil
}

func (l *GameLauncher) startGame(u *url.URL, config Config, launchType LaunchType, cmdBuilder CommandBuilder) error {
	args, err := cmdBuilder.GetArgs(l.repository, u, config, launchType)
	if err != nil {
		return err
	}

	path := filepath.Join(config.InstallPath, config.ExecutablePath, config.ExecutableName)

	dir := config.InstallPath
	if config.StartIn == LaunchDirBinaryDir {
		// Launch in binary directory instead of install path if requested
		dir = filepath.Dir(path)
	}

	cmd := &exec.Cmd{
		Dir:  dir,
		Path: path,
		Args: append([]string{path}, args...),
	}

	return cmd.Start()
}

func toHookHandlerMap(hookHandlers []HookHandler) map[string]HookHandler {
	hookHandlerMap := map[string]HookHandler{}
	for _, h := range hookHandlers {
		hookHandlerMap[h.String()] = h
	}
	return hookHandlerMap
}
