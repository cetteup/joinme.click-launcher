package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
	"net/url"
)

var frostbite3DefaultArgs = []string{
	"-gameMode", "MP",
	"-role", "soldier",
	"-asSpectator", "false",
	"-joinWithParty", "false",
}

var PlusConnectCmdBuilder launcher.CommandBuilder = func(scheme string, host string, port string, u *url.URL) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", host, port)}, nil
}
