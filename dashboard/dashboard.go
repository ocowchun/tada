package dashboard

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
	// tadaConfig "github.com/ocowchun/tada/config"
	"github.com/ocowchun/tada/utils"
	widget "github.com/ocowchun/tada/widget"
	"github.com/ocowchun/tada/widgets/foo"
	"github.com/ocowchun/tada/widgets/ghpr"
	"github.com/ocowchun/tada/widgets/ghreview"
	"github.com/rivo/tview"
)

type Dashboard struct {
	app            *tview.Application
	widgets        []*widget.Widget
	widgetIdx      int
	hasFocusWidget bool
}

func (d *Dashboard) hover(oldViewIdx, newViewIdx int) {
	widgets := d.widgets
	currentFocusedItem := widgets[oldViewIdx]
	currentFocusedItem.Unhover()
	focusedItem := widgets[newViewIdx]
	focusedItem.Hover()
}

func (d *Dashboard) focus() {
	widget := d.widgets[d.widgetIdx]
	d.app.SetFocus(widget)
	d.hasFocusWidget = true
}

func (d *Dashboard) blur() {
	widget := d.widgets[d.widgetIdx]
	widget.Blur()
	d.app.SetFocus(nil)
	d.hasFocusWidget = false
}

type WithBox interface {
	Box() *tview.Box
}

func inputCaptureFactory(d *Dashboard) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if d.hasFocusWidget == false {
			switch event.Key() {
			case tcell.KeyRight:
				oldIdx := d.widgetIdx
				d.widgetIdx = d.widgetIdx + 1
				if d.widgetIdx >= len(d.widgets) {
					d.widgetIdx = 0
				}
				d.hover(oldIdx, d.widgetIdx)
			case tcell.KeyLeft:
				oldIdx := d.widgetIdx
				d.widgetIdx = d.widgetIdx - 1
				if d.widgetIdx < 0 {
					d.widgetIdx = len(d.widgets) - 1
				}
				d.hover(oldIdx, d.widgetIdx)
			case tcell.KeyEnter:
				d.focus()
			}
		}

		if event.Rune() == 'q' {
			if d.hasFocusWidget {
				d.hasFocusWidget = false
				d.blur()
			} else {
				d.app.Stop()
			}
		}
		return event
	}
}

func (d *Dashboard) Run() {
	app := tview.NewApplication()
	d.app = app
	basePath := utils.FindBasePath()
	path := basePath + "/tada.toml"
	config := LoadConfig(path)
	newPrimitive := func(text string) tview.Primitive {
		view := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
		view.SetBorder(true)
		view.SetDynamicColors(true)
		view.SetBorderColor(tcell.ColorDarkRed)
		fmt.Fprintf(view, "%s ", "[red]the[white]")
		return view
	}

	pages := tview.NewPages()

	widget2 := tview.NewGrid().
		SetBorders(true)
	widget2.SetBorderColor(tcell.ColorYellow)
	widget2.AddItem(newPrimitive("widget2 Header"), 0, 0, 4, 3, 0, 0, false)

	grid := tview.NewGrid().
		SetColumns(2, 0, 0, 0, 0, 0, 0, 2).
		SetRows(2, 0, 0, 0, 0, 2).
		SetBorders(false)
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).

	// Layout for screens wider than 100 cells.
	pages.AddPage("grid", grid, true, true)

	widgets := []*widget.Widget{}
	buildinWidgets := map[string]func(config widget.Config, stop func()) *widget.Widget{
		"tada-github":        ghpr.NewWidget,
		"tada-github-review": ghreview.NewWidget,
		"tada-foo":           foo.NewWidget,
	}

	for _, widgetConfig := range config.Widgets {
		newWidget := buildinWidgets[widgetConfig.Name]
		var primitive tview.Primitive
		var w *widget.Widget
		if newWidget != nil {
			w = newWidget(widgetConfig, app.Stop)
			primitive = w
		} else {
			box := LoadPlugin(widgetConfig.Name, widgetConfig)
			w = widget.NewWidget(box)
			primitive = w
		}
		grid.AddItem(primitive, widgetConfig.Y+1, widgetConfig.X+1, widgetConfig.Height,
			widgetConfig.Width, 0, 100, false)
		widgets = append(widgets, w)
	}

	d.widgets = widgets
	d.widgetIdx = 0
	app.SetInputCapture(inputCaptureFactory(d))
	go func() {
		// without below line, the program will broken = =
		// foo.NewWidget().Render()
		for {
			for _, w := range d.widgets {
				if !w.IsRendering() {
					w.Render()
				}
			}
			app.Draw()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	if err := app.SetRoot(pages, true).Run(); err != nil {
		fmt.Println("yoyoyo")
		panic(err)
	}
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}
