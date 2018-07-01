package dashboard

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/ocowchun/tada/widgets"
	"github.com/ocowchun/tada/widgets/foo"
	"github.com/ocowchun/tada/widgets/github"
	"github.com/rivo/tview"
	// "github.com/senorprogrammer/wtf/github"
)

type Dashboard struct {
	app *tview.Application
	// widgets []tview.Primitive
	widgets []*widgets.Widget
	// current hover widget idx
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
	// SetBorder(bool)
	// SetBorderColor(tcell.Color)
}

func inputCaptureFactory(d *Dashboard) func(event *tcell.EventKey) *tcell.EventKey {
	// idx := 0
	// focusedItem := d.widgets[idx]
	// item := focusedItem.(*tview.Grid)
	// // item.Box().SetBorder(true)
	// // item.Box().SetBorderColor(tcell.ColorGreen)

	// item.SetBorder(true)
	// item.SetBorderColor(tcell.ColorGreen)

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
				// d.hover(0, d.widgetIdx)
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
	clockView := tview.NewTextView()
	clockView.SetTextAlign(tview.AlignCenter).SetText("1")

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

	// clockView.SetBorder(true)
	// clockView.SetBorderColor(tcell.ColorGreen)

	pages := tview.NewPages()
	// menu := newPrimitive("Menu")

	widget2 := tview.NewGrid().
		SetBorders(true)
	widget2.SetBorderColor(tcell.ColorYellow)
	widget2.AddItem(newPrimitive("widget2 Header"), 0, 0, 4, 3, 0, 0, false)

	go func() {
		num := 1
		for {
			num = num + 1
			time.Sleep(100 * time.Millisecond)
			clockView.SetText(strconv.Itoa(num))
			app.Draw()
			// app.SetFocus(main)
			// fmt.Println(main.(*tview.TextView).HasFocus())
		}
	}()

	grid := tview.NewGrid().
		SetColumns(2, 0, 0, 0, 0, 0, 0, 2).
		SetRows(2, 0, 0, 0, 0, 2).
		SetBorders(false)
	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	// grid.AddItem(main, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	// grid.AddItem(main, 0, 1, 1, 1, 0, 100, false)
	// grid.AddItem(widget2, 1, 1, 1, 1, 0, 100, false)
	pages.AddPage("grid", grid, true, true)

	// box1 := tview.NewGrid().
	// 	SetRows(0, 0, 0, 0).
	// 	SetColumns(0, 0, 0, 0)
	// box1.SetBorder(true).
	// 	SetBorderColor(tcell.ColorGreen)
	// box1.AddItem(newPrimitive("box1"), 0, 0, 1, 1, 0, 0, false)

	// │,┌,─, ,
	// cells := [][]string{}
	// newHorizontalLine := func(length int) string {
	// 	line := ""
	// 	max := length / 2
	// 	for i := 0; i < max; i++ {
	// 		line = line + "─"
	// 	}
	// 	if length%2 == 1 {
	// 		line = line + "─"
	// 	}
	// 	return line
	// }

	// topLine := newHorizontalLine(28)
	// box1 := tview.NewTextView()
	// box1.SetText("┌" + topLine + "┐")
	box1 := github.NewWidget()
	box2 := foo.NewWidget()
	// box2 := github.NewWidget()
	// box3 := github.NewWidget()
	// box3 := newPrimitive("box3")
	// box4 := newPrimitive("box4")
	// box5 := newPrimitive("box5")

	grid.AddItem(box2, 1, 1, 1, 1, 0, 100, false)
	// grid.AddItem(box3, 2, 1, 1, 1, 0, 100, false)
	// grid.AddItem(box4, 3, 1, 1, 1, 0, 100, false)
	// grid.AddItem(box5, 4, 1, 1, 1, 0, 100, false)

	grid.AddItem(box1, 2, 2, 2, 3, 0, 100, false)

	// pages.
	// gridItems := []tview.Primitive{menu, main, sideBar, clockView}
	// gridItems := []tview.Primitive{main, widget2}
	d.widgets = []*widgets.Widget{box1}
	d.widgetIdx = 0
	app.SetInputCapture(inputCaptureFactory(d))
	// app.SetFocus(main)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		fmt.Println("yoyoyo")
		// fmt.Println(main.GetFocusable())
		panic(err)
	}
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}
