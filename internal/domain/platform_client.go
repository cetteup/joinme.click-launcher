package domain

import (
	"github.com/cetteup/joinme.click-launcher/pkg/game_launcher"
	"github.com/cetteup/joinme.click-launcher/pkg/software_finder"
)

type Platform string

const (
	PlatformOrigin Platform = "Origin"
	PlatformEaApp  Platform = "EA App"
)

type PlatformClient struct {
	Platform       Platform
	FinderConfig   software_finder.Config
	LauncherConfig game_launcher.Config
}
