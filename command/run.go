package command

import (
	"github.com/mitchellh/cli"
	"log"

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
}

func newListWidget() Widget {
	//
	return newList()
}

func (*RunCommand) Run(args []string) int {
	//d := dashboard.NewDashboard()
	//d.Run()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	//events :
	l := newListWidget()

	defaultBorderStyle := ui.NewStyle(ui.ColorWhite)
	//focusedBorderStyle := ui.NewStyle(ui.ColorYellow)

	l.SetBorderStyle(defaultBorderStyle)
	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		switch e.ID {
		case "q", "<C-c>":
			return 0
		}
		l.HandleUIEvent(e)
		ui.Render(l)

	}

	return 0
}

func (*RunCommand) Synopsis() string {
	return "Run tada dashboard"
}

func RunCommandFactory() (cli.Command, error) {
	return &RunCommand{}, nil
}
