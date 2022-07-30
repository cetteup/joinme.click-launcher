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

const (
	LaunchDirInstallDir LaunchDir = iota
	LaunchDirBinaryDir
)

type Config struct {
	DefaultArgs    []string
	StartIn        LaunchDir
	ExecutableName string
	// Relative path from install path to folder containing the executable
	ExecutablePath string
	InstallPath    string
	// Indicates whether game needs to be "cold-launched" (whether we need to kill any pre-existing game instances)
	CloseBeforeLaunch bool
	// Additional process names related to this game (some games start multiple processes, which we need to kill pre-launch)
	AdditionalProcessNames map[string]bool
}

type URLValidator func(u *url.URL) error

type CommandBuilder func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error)

func PrepareLaunch(config Config) error {
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

func StartGame(config Config, cmdBuilder CommandBuilder, scheme string, host string, port string, u *url.URL) error {
	args, err := cmdBuilder(config.InstallPath, scheme, host, port, u)
	if err != nil {
		return err
	}

	args = append(args, config.DefaultArgs...)

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
