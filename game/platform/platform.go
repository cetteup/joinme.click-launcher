package platform

import (
	"github.com/cetteup/joinme.click-launcher/game/finder"
	"github.com/cetteup/joinme.click-launcher/game/launcher"
)

type Platform string

const (
	Origin Platform = "Origin"
)

type Client struct {
	Platform       Platform
	FinderConfig   finder.Config
	LauncherConfig launcher.Config
}
