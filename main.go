package main

import (
	"fmt"
	"github.com/mefellows/godspeed/command"
	_ "github.com/mefellows/godspeed/deployment/strategy"
	_ "github.com/mefellows/godspeed/repository"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

func main() {
	cli := cli.NewCLI(strings.ToLower(APPLICATION_NAME), VERSION)
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
