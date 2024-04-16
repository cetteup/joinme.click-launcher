package internal

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

const (
	HookKillProcess         = "kill-process"
	HookDeleteFile          = "delete-running-file" // CoD games write a dummy file when launched. If the file is still present when launched again, the game assumes it crashed and offers to start in safe mode.
	PlusConnectPrefix       = "+connect"
	Frostbite3GameIdPattern = `^\d+$` // game ids vary by length, so for now we are just validating that it only contains numbers
)

type IPPortURLValidator struct{}

func (v IPPortURLValidator) Validate(u *url.URL) error {
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

func MakePatternURLValidator(pattern string) PatternURLValidator {
	return PatternURLValidator{
		pattern: pattern,
	}
}

type PatternURLValidator struct {
	pattern string
}

func (v PatternURLValidator) Validate(u *url.URL) error {
	hostname := u.Hostname()
	matched, err := regexp.Match(v.pattern, []byte(hostname))
	if err != nil {
		return fmt.Errorf("failed to validate game id: %s", err)
	}
	if !matched {
		return fmt.Errorf("url hostname is not a valid game id: %s", hostname)
	}

	return nil
}

func MakeSimpleCmdBuilder(prefixes ...string) SimpleCmdBuilder {
	return SimpleCmdBuilder{
		prefixes: prefixes,
	}
}

type SimpleCmdBuilder struct {
	prefixes []string
}

func (b SimpleCmdBuilder) GetArgs(_ game_launcher.FileRepository, u *url.URL, launchType game_launcher.LaunchType) ([]string, error) {
	args := make([]string, 0, len(b.prefixes))
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args, b.prefixes...)
		args = append(args, net.JoinHostPort(u.Hostname(), u.Port()))
	}

	return args, nil
}

type OriginCmdBuilder struct {
}

func (b OriginCmdBuilder) GetArgs(_ game_launcher.FileRepository, u *url.URL, launchType game_launcher.LaunchType) ([]string, error) {
	args := make([]string, 0, 8)
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args,
			"-gameMode", "MP",
			"-role", "soldier",
			"-asSpectator", "false",
			"-gameId", u.Hostname(),
		)
	}

	return args, nil
}

type RefractorV1CmdBuilder struct{}

func (b RefractorV1CmdBuilder) GetArgs(_ game_launcher.FileRepository, u *url.URL, launchType game_launcher.LaunchType) ([]string, error) {
	args := make([]string, 0, 6)
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args, "+joinServer", u.Hostname(), "+port", u.Port())
	}

	query := u.Query()
	if internal.QueryHasMod(query) {
		args = append(args, "+game", internal.GetModFromQuery(query))
	}

	return args, nil
}

// MakeKillProcessHookHandler Returns a hook handler that kills any running game processes plus any additional targets
func MakeKillProcessHookHandler(targetLaunchExecutable bool, additionalTargets ...string) KillProcessHookHandler {
	return KillProcessHookHandler{
		targetLaunchExecutable: targetLaunchExecutable,
		additionalTargets:      additionalTargets,
	}
}

type KillProcessHookHandler struct {
	targetLaunchExecutable bool
	additionalTargets      []string
}

func (h KillProcessHookHandler) Run(_ game_launcher.FileRepository, _ *url.URL, config game_launcher.Config, _ game_launcher.LaunchType, _ map[string]string) error {
	targets := h.additionalTargets
	if h.targetLaunchExecutable {
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

func (h KillProcessHookHandler) String() string {
	return HookKillProcess
}

func MakeDeleteFileHookHandler(pathsBuilder func(config game_launcher.Config) ([]string, error)) DeleteFileHookHandler {
	return DeleteFileHookHandler{
		pathsBuilder: pathsBuilder,
	}
}

type DeleteFileHookHandler struct {
	pathsBuilder func(config game_launcher.Config) ([]string, error)
}

func (h DeleteFileHookHandler) Run(fr game_launcher.FileRepository, _ *url.URL, config game_launcher.Config, _ game_launcher.LaunchType, _ map[string]string) error {
	paths, err := h.pathsBuilder(config)
	if err != nil {
		return err
	}
	for _, path := range paths {
		if err := DeleteFileIfExists(fr, path); err != nil {
			return err
		}
	}
	return nil
}

func (h DeleteFileHookHandler) String() string {
	return HookDeleteFile
}

func CoDRunningFilePathsBuilder(config game_launcher.Config) ([]string, error) {
	// Running file name will be the executable name minus the .exe prefixed by __, e.g. "__CoDMP"
	name := fmt.Sprintf("__%s", strings.TrimSuffix(config.ExecutableName, filepath.Ext(config.ExecutableName)))

	// Primary place to look is in the install path, basically right "next to" the executable
	primary := filepath.Join(config.InstallPath, name)

	// At least when installed in the default location, CoD2 may store the file in the VirtualStore
	virtualStore, err := buildVirtualStorePath()
	if err != nil {
		return nil, err
	}
	// With the game installed in C:\Program Files\Call of Duty, the alternate would be ...\AppData\Local\VirtualStore\Program Files\Call of Duty
	alternate := filepath.Join(
		virtualStore,
		strings.TrimPrefix(config.InstallPath, filepath.VolumeName(config.InstallPath)),
		name,
	)

	return []string{primary, alternate}, nil
}
