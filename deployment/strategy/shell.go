package strategy

import (
	"errors"
	"fmt"
	"github.com/mefellows/godspeed/log"
	"github.com/mefellows/plugo/plugo"
	"math/rand"
	"os/exec"
	"time"
)

type ShellDeploymentStrategy struct {
	Host             string `mapstructure:"ssh_host"`
	Shell            string
	ShellArgs        []string `mapstructure:"shell_args"`
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

func (s ShellDeploymentStrategy) createCommand(command string) *exec.Cmd {
	// Create command string
	args := append(s.ShellArgs, fmt.Sprintf(`"%s"`, command))
	return exec.Command(s.Shell, args...)
}

func (s ShellDeploymentStrategy) runCommand(c string) error {

	// Create command string
	cmd := s.createCommand(c)
	log.Info(" --> Running command: %s", c)

	// Create command string
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Script '%s' with args '%v' not found or is not executable: %v", cmd.Path, cmd.Args, err))
		return err
	}
	log.Info(log.Colorize(log.CYAN, fmt.Sprintf("%s", out)))
	return nil
}

func (s ShellDeploymentStrategy) Deploy() error {
	log.Info("Shell Deploying")
	for _, c := range s.Commands {
		s.runCommand(c)
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
		s.runCommand(c)
	}
	log.Info("Shell Rollback complete!")
	return nil
}

func (s ShellDeploymentStrategy) Teardown() {
	log.Debug("Shell Teardown")
}
