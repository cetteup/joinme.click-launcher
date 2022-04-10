package game

type Router struct {
	repository RegistryRepository
	Launchers  map[string]*Launcher
}

func NewRouter(repository RegistryRepository) *Router {
	return &Router{
		repository: repository,
		Launchers:  map[string]*Launcher{},
	}
}

func (r Router) AddLauncher(config LauncherConfig, cmdBuilder CommandBuilder) {
	launcher := NewLauncher(r.repository, config, cmdBuilder)
	r.Launchers[launcher.Config.ProtocolScheme] = launcher
}
