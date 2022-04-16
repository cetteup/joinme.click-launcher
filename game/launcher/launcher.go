package launcher

import (
	"net/url"
	"os/exec"
	"path/filepath"
)

type LaunchDir int

const (
	InstallDir LaunchDir = iota
	BinaryDir
)

type Config struct {
	DefaultArgs    []string
	StartIn        LaunchDir
	ExecutablePath string
	InstallPath    string
}

type CommandBuilder func(scheme string, host string, port string, u *url.URL) ([]string, error)

type GameLauncher struct {
	Config     Config
	CmdBuilder CommandBuilder
}

func NewGameLauncher(config Config, cmdBuilder CommandBuilder) *GameLauncher {
	return &GameLauncher{
		Config:     config,
		CmdBuilder: cmdBuilder,
	}
}

func (l *GameLauncher) StartGame(scheme string, host string, port string, u *url.URL) error {
	args, err := l.CmdBuilder(scheme, host, port, u)
	if err != nil {
		return err
	}

	args = append(args, l.Config.DefaultArgs...)

	path := filepath.Join(l.Config.InstallPath, l.Config.ExecutablePath)

	dir := l.Config.InstallPath
	if l.Config.StartIn == BinaryDir {
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
