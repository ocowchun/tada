package command

import (
	"github.com/mitchellh/cli"
	"github.com/ocowchun/tada/dashboard"
)

type RunCommand struct {
}

func (*RunCommand) Help() string {
	return "Run tada dashboard, you must run `tada init` to create config first."
}

func (*RunCommand) Run(args []string) int {
	//d := dashboard.New()
	//d.Run()
	d := dashboard.New()
	err := d.Run()

	if err != nil {
		return 1
	}
	return 0
}

func (*RunCommand) Synopsis() string {
	return "Run tada dashboard"
}

func RunCommandFactory() (cli.Command, error) {
	return &RunCommand{}, nil
}
