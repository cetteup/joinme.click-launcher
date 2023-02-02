package internal

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows"

	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

const (
	virtualStoreDirName = "VirtualStore"
)

func buildOriginURL(offerIDs []string, args []string) string {
	params := url.Values{
		"offerIds":  {strings.Join(offerIDs, ",")},
		"authCode":  {},
		"cmdParams": {url.PathEscape(strings.Join(args, " "))},
	}
	u := url.URL{
		Scheme:   "origin2",
		Path:     "game/launch",
		RawQuery: params.Encode(),
	}
	return u.String()
}

func isTargetProcess(targets []string, executable string) bool {
	for _, target := range targets {
		if executable == target {
			return true
		}
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

func DeleteFileIfExists(fr game_launcher.FileRepository, path string) error {
	// Make sure it's a file, so we don't accidentally delete something else
	exists, err := fr.FileExists(path)
	if err != nil {
		return err
	}
	if exists {
		return fr.RemoveAll(path)
	}
	return nil
}

func GetLocalAppDataPath() (string, error) {
	return windows.KnownFolderPath(windows.FOLDERID_LocalAppData, windows.KF_FLAG_DEFAULT)
}

func buildVirtualStorePath() (string, error) {
	appData, err := GetLocalAppDataPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(appData, virtualStoreDirName), nil
}
