package command

import (
	"github.com/mitchellh/cli"
	//"github.com/ocowchun/tada/dashboard"
	ui "github.com/gizak/termui/v3"
)

type RunCommand struct {
}

func (*RunCommand) Help() string {
	return "Run tada dashboard, you must run `tada init` to create config first."
}

type Widget interface {
	HandleUIEvent(event ui.Event)
	ui.Drawable
	SetBorderStyle(style ui.Style)
	//SetRect(x1, y1, x2, y2 int)
}

func newListWidget() Widget {
	//
	return newList()
}

func (*RunCommand) Run(args []string) int {
	//d := dashboard.NewDashboard()
	//d.Run()
	d := NewDashboard()
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
