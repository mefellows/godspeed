package command

import (
	"flag"
	"github.com/mefellows/godspeed/godspeed"
	"github.com/mefellows/godspeed/log"
	"strings"
)

type DeployCommand struct {
	Meta Meta
}

func (pc *DeployCommand) Run(args []string) int {
	cmdFlags := flag.NewFlagSet("proxy", flag.ContinueOnError)
	cmdFlags.Usage = func() { pc.Meta.Ui.Output(pc.Help()) }
	c := &godspeed.GodspeedConfig{}

	cmdFlags.StringVar(&c.ConfigFile, "config", "", "Path to a YAML configuration file")

	// Validate
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	godspeed := godspeed.New(c)
	godspeed.Setup()

	log.Info("Performing Deployment...")

	for _, d := range godspeed.DeploymentStrategies {
		res := d.Deploy()
		if res != nil {
			d.Rollback()
		}
	}

	godspeed.Shutdown()

	return 0
}

func (c *DeployCommand) Help() string {
	helpText := `
Usage: godspeed deploy [options] 

  Deploy application
  
Options:

  --config                    Location of Godspeed configuration file
`

	return strings.TrimSpace(helpText)
}

func (c *DeployCommand) Synopsis() string {
	return "Run the Godspeed deployer"
}
