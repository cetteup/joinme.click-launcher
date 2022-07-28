package titles

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

const (
	// game ids vary by length, so for now we are just validating that it only contains numbers
	frostbite3GameIdPattern = `^\d+$`
	urlQueryKeyMod          = "mod"
)

var ipPortURLValidator launcher.URLValidator = func(u *url.URL) error {
	hostname, port := u.Hostname(), u.Port()
	if !isValidIPv4(hostname) {
		return fmt.Errorf("url hostname is not a valid IPv4 address: %s", hostname)
	}
	if port == "" {
		return fmt.Errorf("port is missing from url")
	}
	// When parsing a URL, only the port format is validated (numbers only)
	// The url package does not ensure that a port is within the valid TCP/UDP port range, so we need to take care of that
	if !isValidPort(port) {
		return fmt.Errorf("url port is not a valid network port: %s", port)
	}

	return nil
}

var frostbite3GameIdURLValidator launcher.URLValidator = func(u *url.URL) error {
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

var frostbite3DefaultArgs = []string{
	"-gameMode", "MP",
	"-role", "soldier",
	"-asSpectator", "false",
	"-joinWithParty", "false",
}

var plusConnectCmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", host, port)}, nil
}
