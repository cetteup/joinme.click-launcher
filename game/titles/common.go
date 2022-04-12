package titles

import (
	"fmt"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

var PlusConnectCmdBuilder launcher.CommandBuilder = func(scheme string, ip string, port string) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", ip, port)}, nil
}
