package game_launcher

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"
)

type LaunchDir int
type LaunchType int

const (
	LaunchDirInstallDir LaunchDir = iota
	LaunchDirBinaryDir

	LaunchTypeLaunchAndJoin LaunchType = iota
	LaunchTypeLaunchOnly
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
	// Indicates whether game needs to be "cold-launched" (whether we need to kill any pre-existing game instances)
	CloseBeforeLaunch bool
	// Additional process names related to this game (some games start multiple processes, which we need to kill pre-launch)
	AdditionalProcessNames map[string]bool
}

type URLValidator func(u *url.URL) error

// CommandBuilder Function to construct slice of launch arguments for a game. Receives a file repository in order to be
// able to access any config file it may need to construct the arguments.
type CommandBuilder func(fr FileRepository, u *url.URL, config Config, launchType LaunchType) ([]string, error)

func (l *GameLauncher) PrepareLaunch(config Config) error {
	if !config.CloseBeforeLaunch {
		return nil
	}

	processes, err := ps.Processes()
	if err != nil {
		return fmt.Errorf("failed to retrieve process list: %s", err)
	}

	killed := map[int]string{}
	for _, process := range processes {
		if isGameProcess(config, process) {
			log.Info().
				Int("pid", process.Pid()).
				Str("executable", process.Executable()).
				Msg("Killing existing game process")
			if err = killProcess(process.Pid()); err != nil {
				return fmt.Errorf("failed to kill existing game process %s (%d): %s", process.Executable(), process.Pid(), err)
			}
			killed[process.Pid()] = process.Executable()
		}
	}

	// Wait for killed processes to exit
	if err = waitForProcessesToExit(killed); err != nil {
		return err
	}

	return nil
}

func (l *GameLauncher) StartGame(u *url.URL, config Config, launchType LaunchType, cmdBuilder CommandBuilder) error {
	args, err := cmdBuilder(l.repository, u, config, launchType)
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

func isGameProcess(config Config, process ps.Process) bool {
	if process.Executable() == config.ExecutableName {
		return true
	}
	if _, isAdditionalGameProcess := config.AdditionalProcessNames[process.Executable()]; isAdditionalGameProcess {
		return true
	}

	return false
}

func killProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err = proc.Signal(syscall.SIGKILL); err != nil {
		return err
	}

	return nil
}

func waitForProcessesToExit(processes map[int]string) error {
	iterations := 0
	for ; len(processes) > 0 && iterations < 5; iterations++ {
		for pid, executable := range processes {
			log.Debug().
				Int("pid", pid).
				Str("executable", executable).
				Msg("Checking if game process exited")
			proc, err := ps.FindProcess(pid)
			if err != nil {
				return fmt.Errorf("failed to check if killed game process is still running: %s", err)
			}

			// Remove process from map if it exited (was no longer found)
			if proc == nil {
				log.Debug().
					Int("pid", pid).
					Str("executable", executable).
					Msg("Game process is gone")
				delete(processes, pid)
			}
		}
		time.Sleep(1 * time.Second)
	}

	// Return error if not all processes exited yet
	if len(processes) > 0 {
		return fmt.Errorf("timed out waiting for killed game processes to exit")
	}

	return nil
}
