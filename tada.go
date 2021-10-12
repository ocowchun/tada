package main

import (
	"github.com/mitchellh/cli"
	"github.com/ocowchun/tada/command"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("app", "0.2.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"run": command.RunCommandFactory,
		//"init":    command.InitCommandFactory,
	}
	//
	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
