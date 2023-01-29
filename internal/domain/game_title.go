package domain

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cetteup/joinme.click-launcher/internal"
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

type GameTitle struct {
	Name           string
	ProtocolScheme string
	PlatformClient *PlatformClient
	FinderConfigs  []software_finder.Config
	Mods           []GameMod
	LauncherConfig game_launcher.Config
	URLValidator   game_launcher.URLValidator
	CmdBuilder     game_launcher.CommandBuilder
	HookHandlers   []game_launcher.HookHandler
}

func (t *GameTitle) AddCustomConfig(config internal.CustomLauncherConfig) {
	if config.HasExecutableName() {
		t.LauncherConfig.ExecutableName = config.ExecutableName
	}

	if config.HasExecutablePath() {
		t.LauncherConfig.ExecutablePath = config.ExecutablePath
	}

	if config.HasInstallPath() {
		// Prepend custom path based finder in order search any custom paths first
		t.FinderConfigs = append([]software_finder.Config{
			{
				ForType:     software_finder.PathFinder,
				InstallPath: config.InstallPath,
				PathType:    software_finder.PathTypeDir,
			},
		}, t.FinderConfigs...)
	}

	if config.HasArgs() {
		t.LauncherConfig.DefaultArgs = append(t.LauncherConfig.DefaultArgs, config.Args...)
	}

	if config.HasHookConfigs() {
		for _, hook := range config.Hooks {
			t.LauncherConfig.HookConfigs = append(t.LauncherConfig.HookConfigs, game_launcher.HookConfig{
				Handler:     hook.Handler,
				When:        hook.When,
				ExitOnError: hook.ExitOnError,
				Args:        hook.Args,
			})
		}
	}
}

func (t *GameTitle) RequiresPlatformClient() bool {
	return t.PlatformClient != nil
}

func (t *GameTitle) GetMod(slug string) *GameMod {
	for _, mod := range t.Mods {
		if strings.EqualFold(slug, mod.Slug) {
			return &mod
		}
	}
	return nil
}

func (t *GameTitle) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("%s (%s)", t.Name, t.ProtocolScheme)
}

func MakeMod(name string, slug string, finderConfigs []software_finder.Config) GameMod {
	return GameMod{
		Name:          name,
		Slug:          slug,
		finderConfigs: finderConfigs,
	}
}

type GameMod struct {
	Name          string
	Slug          string
	finderConfigs []software_finder.Config
}

// ComputeFinderConfigs Mod finder configs (can) only contain relative paths based on the game's install dir.
// So, we need to compute absolute paths for any software_finder.PathFinder configs before we can use them.
func (m *GameMod) ComputeFinderConfigs(gameInstallPath string) []software_finder.Config {
	computedConfigs := make([]software_finder.Config, 0, len(m.finderConfigs))
	for _, config := range m.finderConfigs {
		// Config is not a pointer, so we can change "it" and the function call remains idempotent
		if config.ForType == software_finder.PathFinder {
			config.InstallPath = filepath.Join(gameInstallPath, config.InstallPath)
		}
		computedConfigs = append(computedConfigs, config)
	}
	return computedConfigs
}

func (m *GameMod) String() string {
	if m == nil {
		return "nil"
	}
	return fmt.Sprintf("%s (%s)", m.Name, strings.ToLower(m.Slug))
}
