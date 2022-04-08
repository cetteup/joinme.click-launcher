package game

import (
	"fmt"
)

var PlusConnectCmdBuilder CommandBuilder = func(config LauncherConfig, ip string, port string) ([]string, error) {
	return []string{"+connect", fmt.Sprintf("%s:%s", ip, port)}, nil
}
