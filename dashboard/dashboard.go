package dashboard

import (
	ui "github.com/gizak/termui/v3"
	widget "github.com/ocowchun/tada/Widget"
	widgets "github.com/ocowchun/tada/widgets"
	"github.com/ocowchun/tada/widgets/github_pr_list"
	"log"
)

type Dashboard struct {
}

var (
	defaultBorderStyle  = ui.NewStyle(ui.ColorWhite)
	selectedBorderStyle = ui.NewStyle(ui.ColorYellow)
	focusedBorderStyle  = ui.NewStyle(ui.ColorGreen)
)

func (d Dashboard) Run() error {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// TODO: figure out a better way to init widgets
	l := github_pr_list.NewGitHubPRList()
	l.SetRect(0, 0, 55, 8)
	l.SetBorderStyle(defaultBorderStyle)

	l2 := widgets.NewList()
	l2.SetRect(10, 10, 35, 18)
	l2.SetBorderStyle(defaultBorderStyle)

	widgets := []widget.Widget{l, l2}

	ds := make([]ui.Drawable, 0, len(widgets))
	for _, w := range widgets {
		ds = append(ds, w)
	}
	//ds = append(ds, widgets[0])
	//ds = append(ds, widgets[1])

	ui.Render(ds...)

	uiEvents := ui.PollEvents()

	var currentWidget widget.Widget
	currentWidgetIdx := -1
	widgetMode := false
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Escape>":
			return nil
		}

		if !widgetMode {
			if currentWidget != nil {
				currentWidget.SetBorderStyle(defaultBorderStyle)
			}
			switch e.ID {
			case "<Right>":
				currentWidgetIdx++
				if currentWidgetIdx == len(widgets) {
					currentWidgetIdx = 0
				}
				currentWidget = widgets[currentWidgetIdx]
				currentWidget.SetBorderStyle(selectedBorderStyle)
			case "<Left>":
				currentWidgetIdx--
				if currentWidgetIdx < 0 {
					currentWidgetIdx = len(widgets) - 1
				}
				currentWidget = widgets[currentWidgetIdx]
				currentWidget.SetBorderStyle(selectedBorderStyle)
			case "<Enter>":
				if currentWidget != nil {
					widgetMode = true
					currentWidget.SetBorderStyle(focusedBorderStyle)
				}
			case "r":
				//	TODO: refresh all components
			}

		} else {
			if e.ID == "q" {
				currentWidget.SetBorderStyle(selectedBorderStyle)
				widgetMode = false
			} else {
				currentWidget.HandleUIEvent(e)
			}
		}
		ui.Render(ds...)

	}
	return nil
}

func New() *Dashboard {
	return &Dashboard{}
}
