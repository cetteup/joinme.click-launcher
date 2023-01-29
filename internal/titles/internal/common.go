package internal

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
)

const (
	HookKillProcess         = "kill-process"
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

func (b SimpleCmdBuilder) GetArgs(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args := append(b.prefixes, net.JoinHostPort(u.Hostname(), u.Port()))
		if config.AppendDefaultArgs {
			return append(args, config.DefaultArgs...), nil
		}
		return append(config.DefaultArgs, args...), nil
	}
	return nil, nil
}

func MakeOriginCmdBuilder(offerIDs ...string) OriginCmdBuilder {
	return OriginCmdBuilder{
		offerIDs: offerIDs,
	}
}

type OriginCmdBuilder struct {
	offerIDs []string
}

func (b OriginCmdBuilder) GetArgs(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	args := config.DefaultArgs
	if launchType == game_launcher.LaunchTypeLaunchAndJoin {
		args = append(args,
			"-gameMode", "MP",
			"-role", "soldier",
			"-asSpectator", "false",
			"-joinWithParty", "false",
			"-gameId", u.Hostname(),
		)
	}

	originURL := buildOriginURL(b.offerIDs, args)
	return []string{originURL}, nil
}

type RefractorV1CmdBuilder struct{}

func (b RefractorV1CmdBuilder) GetArgs(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType) ([]string, error) {
	args := config.DefaultArgs
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

func (h KillProcessHookHandler) Run(fr game_launcher.FileRepository, u *url.URL, config game_launcher.Config, launchType game_launcher.LaunchType, args map[string]string) error {
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
