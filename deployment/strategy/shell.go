package strategy

import (
	"errors"
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
	"math/rand"
	"time"
)

type ShellDeploymentStrategy struct {
	Host             string `mapstructure:"ssh_host"`
	Commands         []string
	RollbackCommands []string `mapstructure:"rollback_commands"`
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	plugo.PluginFactories.Register(func() (interface{}, error) {
		return &ShellDeploymentStrategy{}, nil
	}, "shell")
}

func (s ShellDeploymentStrategy) Setup() {
	log.Debug("Setting up Shell")
}

func (s ShellDeploymentStrategy) Deploy() error {
	log.Info("Shell Deploying")
	for _, c := range s.Commands {
		log.Info(" --> Running command: %s", c)
	}

	if rand.Intn(2) == 1 {
		log.Info("Shell Deployment complete with errors, rolling back!")
		return errors.New("Shit, something went wrong!")
	} else {
		log.Info("Shell Deployment complete!")
		return nil
	}
}

func (s ShellDeploymentStrategy) Rollback() error {
	log.Info("Shell Rolling back")
	for _, c := range s.RollbackCommands {
		log.Info(" --> Running command: %s", c)
	}
	log.Info("Shell Rollback complete!")
	return nil
}

func (s ShellDeploymentStrategy) Teardown() {
	log.Debug("Shell Teardown")
}
