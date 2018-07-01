package command

import (
	"github.com/mitchellh/cli"
	"github.com/ocowchun/tada/dashboard"
)

type FooCommand struct {
}

func (*FooCommand) Help() string {
	return "Sorry, there no help"
}

func (*FooCommand) Run(args []string) int {
	d := dashboard.NewDashboard()
	d.Run()

	return 0
}

func (*FooCommand) Synopsis() string {
	return "Run tada dashboard"
}

func RunCommandFactory() (cli.Command, error) {
	return &FooCommand{}, nil
}
