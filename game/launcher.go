package game

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"path/filepath"
)

type LaunchDir int

const (
	RegPathSoftware   = "SOFTWARE"
	RegPathClasses    = "Classes"
	RegPathOpen       = "open"
	RegPathShell      = "shell"
	RegPathCommand    = "command"
	RegKeyDefault     = ""
	RegKeyURLProtocol = "URL Protocol"
)

const (
	InstallDir LaunchDir = iota
	BinaryDir
)

type RegistryRepository interface {
	GetStringValue(k registry.Key, path string, valueName string) (string, error)
	SetStringValue(k registry.Key, path string, valueName string, value string) error
	CreateKey(k registry.Key, path string) error
}

type LauncherConfig struct {
	ProtocolScheme    string
	GameLabel         string
	ExecutablePath    string
	StartIn           LaunchDir
	RegistryPath      string
	RegistryValueName string
}

type CommandBuilder func(config LauncherConfig, ip string, port string) ([]string, error)

type Launcher struct {
	repository RegistryRepository
	Config     LauncherConfig
	CmdBuilder CommandBuilder
}

func NewLauncher(repository RegistryRepository, config LauncherConfig, cmdBuilder CommandBuilder) *Launcher {
	return &Launcher{
		repository: repository,
		Config:     config,
		CmdBuilder: cmdBuilder,
	}
}

func (l *Launcher) IsGameInstalled() (bool, error) {
	_, err := l.repository.GetStringValue(registry.LOCAL_MACHINE, l.Config.RegistryPath, l.Config.RegistryValueName)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (l Launcher) IsHandlerRegistered() (bool, error) {
	path := l.getUrlHandlerRegistryPath([]string{RegPathShell, RegPathOpen, RegPathCommand})
	value, err := l.repository.GetStringValue(registry.CURRENT_USER, path, RegKeyDefault)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	expected, err := l.getHandlerCommand()
	if err != nil {
		return false, err
	}

	return value == expected, nil
}

func (l Launcher) RegisterHandler() error {
	basePath := l.getUrlHandlerRegistryPath(nil)
	err := l.repository.CreateKey(registry.CURRENT_USER, basePath)
	if err != nil {
		return err
	}

	err = l.repository.SetStringValue(registry.CURRENT_USER, basePath, RegKeyDefault, fmt.Sprintf("URL:%s protocol", l.Config.GameLabel))
	if err != nil {
		return err
	}
	err = l.repository.SetStringValue(registry.CURRENT_USER, basePath, RegKeyURLProtocol, "")
	if err != nil {
		return err
	}

	subKeys := []string{RegPathShell, RegPathOpen, RegPathCommand}
	for i := range subKeys {
		subPath := l.getUrlHandlerRegistryPath(subKeys[:i+1])
		err = l.repository.CreateKey(registry.CURRENT_USER, subPath)
		if err != nil {
			return err
		}
	}

	cmdPath := l.getUrlHandlerRegistryPath(subKeys)
	cmd, err := l.getHandlerCommand()
	if err != nil {
		return err
	}

	return l.repository.SetStringValue(registry.CURRENT_USER, cmdPath, RegKeyDefault, cmd)
}

func (l *Launcher) StartGame(ip string, port string) error {
	args, err := l.CmdBuilder(l.Config, ip, port)
	if err != nil {
		return err
	}

	dir, err := l.repository.GetStringValue(registry.LOCAL_MACHINE, l.Config.RegistryPath, l.Config.RegistryValueName)
	if err != nil {
		return err
	}

	path := filepath.Join(dir, l.Config.ExecutablePath)

	// Launch in binary directory instead of install path if requested
	if l.Config.StartIn == BinaryDir {
		dir = filepath.Dir(path)
	}

	cmd := &exec.Cmd{
		Dir:  dir,
		Path: path,
		Args: append([]string{path}, args...),
	}

	return cmd.Start()
}

func (l Launcher) getUrlHandlerRegistryPath(children []string) string {
	path := filepath.Join(RegPathSoftware, RegPathClasses, l.Config.ProtocolScheme)
	for _, child := range children {
		path = filepath.Join(path, child)
	}

	return path
}

func (l Launcher) getHandlerCommand() (string, error) {
	launcherPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\" \"%%1\"", launcherPath), nil
}
