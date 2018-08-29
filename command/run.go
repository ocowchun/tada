package command

import (
	"github.com/mitchellh/cli"
	"github.com/ocowchun/tada/dashboard"
)

type FooCommand struct {
}

func (*FooCommand) Help() string {
	return "Run tada dashboard, you must run `tada init` to create config first."
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
