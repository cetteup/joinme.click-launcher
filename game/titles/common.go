package titles

import (
	"fmt"
	"net/url"

	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

const (
	UrlQueryKeyMod = "mod"
)

var frostbite3DefaultArgs = []string{
	"-gameMode", "MP",
	"-role", "soldier",
	"-asSpectator", "false",
	"-joinWithParty", "false",
}

var plusConnectCmdBuilder launcher.CommandBuilder = func(installPath string, scheme string, host string, port string, u *url.URL) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", host, port)}, nil
}
