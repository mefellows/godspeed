package godspeed

import (
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
)

type GodspeedConfig struct {
	RawConfig  *plugo.RawConfig
	ConfigFile string
}
type PluginConfig struct {
	Name        string
	Description string
	LogLevel    int                  `default:"2" required:"true" mapstructure:"loglevel"`
	Deployment  []plugo.PluginConfig `mapstructure:"deployment"`
}

type Godspeed struct {
	config               *GodspeedConfig
	DeploymentStrategies []DeploymentStrategy
	// Other plugins/abstractions...
	// repositories []Repository
}

func New(config *GodspeedConfig) *Godspeed {
	return &Godspeed{config: config}
}

func NewWithDefaultGodspeedConfig() *Godspeed {
	c := &GodspeedConfig{}
	return &Godspeed{config: c}
}

func (g *Godspeed) Setup() {
	g.LoadPlugins()

	// Setup all plugins...
	for _, p := range g.DeploymentStrategies {
		p.Setup()
	}
}

func (g *Godspeed) Shutdown() {
	for _, p := range g.DeploymentStrategies {
		p.Teardown()
	}
}

func (g *Godspeed) LoadPlugins() {
	// Load Configuration
	var err error
	var confLoader *plugo.ConfigLoader
	c := &PluginConfig{}
	if g.config.ConfigFile != "" {
		confLoader = &plugo.ConfigLoader{}
		err = confLoader.LoadFromFile(g.config.ConfigFile, &c)
		if err != nil {
			log.Fatalf("Unable to read configuration file: %s", err.Error())
		}
	} else {
		log.Fatal("No config file provided")
	}

	log.SetLevel(log.LogLevel(c.LogLevel))

	// Load all plugins
	g.DeploymentStrategies = make([]DeploymentStrategy, len(c.Deployment))
	plugins := plugo.LoadPluginsWithConfig(confLoader, c.Deployment)
	for i, p := range plugins {
		log.Debug("Loading plugin\t" + log.Colorize(log.YELLOW, c.Deployment[i].Name))
		g.DeploymentStrategies[i] = p.(DeploymentStrategy)
	}
}
