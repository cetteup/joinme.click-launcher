package internal

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

const (
	HookKillProcess = "kill-process"
	// game ids vary by length, so for now we are just validating that it only contains numbers
	frostbite3GameIdPattern = `^\d+$`
)

var IPPortURLValidator game_launcher.URLValidator = func(u *url.URL) error {
	hostname, port := u.Hostname(), u.Port()
	if !internal.IsValidIPv4(hostname) {
		return fmt.Errorf("url hostname is not a valid IPv4 address: %s", hostname)
	}
	if port == "" {
		return fmt.Errorf("port is missing from url")
	}
	// When parsing a URL, only the port format is validated (numbers only)
	// The url package does not ensure that a port is within the valid TCP/UDP port range, so we need to take care of that
	if !internal.IsValidPort(port) {
		return fmt.Errorf("url port is not a valid network port: %s", port)
	}

	return nil
}

var Frostbite3GameIdURLValidator game_launcher.URLValidator = func(u *url.URL) error {
	hostname := u.Hostname()
	matched, err := regexp.Match(frostbite3GameIdPattern, []byte(hostname))
	if err != nil {
		return fmt.Errorf("failed to validate game id: %s", err)
	}
	if !matched {
		return fmt.Errorf("url hostname is not a valid game id: %s", hostname)
	}

	return nil
}

var Frostbite3DefaultArgs = []string{
	"-gameMode", "MP",
	"-role", "soldier",
	"-asSpectator", "false",
	"-joinWithParty", "false",
}

var PlusConnectCmdBuilder game_launcher.CommandBuilder = func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args := []string{"+connect", fmt.Sprintf("%s:%s", u.Hostname(), u.Port())}
		if config.AppendDefaultArgs {
			return append(args, config.DefaultArgs...), nil
		}
		return append(config.DefaultArgs, args...), nil
	}
	return nil, nil
}

var PlainCmdBuilder game_launcher.CommandBuilder = func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args := []string{fmt.Sprintf("%s:%s", u.Hostname(), u.Port())}
		if config.AppendDefaultArgs {
			return append(args, config.DefaultArgs...), nil
		}
		return append(config.DefaultArgs, args...), nil
	}
	return nil, nil
}

// KillProcessHookHandler Returns a hook handler that kills any running game processes plus any additional targets
func KillProcessHookHandler(targetLaunchExecutable bool, targets ...string) game_launcher.HookHandler {
	return func(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType, args map[string]string) error {
		if targetLaunchExecutable {
			targets = append(targets, config.ExecutableName)
		}

		processes, err := ps.Processes()
		if err != nil {
			return fmt.Errorf("failed to retrieve process list: %s", err)
		}

		killed := map[int]string{}
		for _, process := range processes {
			if isTargetProcess(targets, process.Executable()) {
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
}
