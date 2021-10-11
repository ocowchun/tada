package widgets

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	widget "github.com/ocowchun/tada/Widget"
)

type MyList struct {
	previousKey string
	widgets.List
}

func (l *MyList) HandleUIEvent(e ui.Event) {
	switch e.ID {
	case "j", "<Down>":
		l.ScrollDown()
	case "k", "<Up>":
		l.ScrollUp()
	case "<C-d>":
		l.ScrollHalfPageDown()
	case "<C-u>":
		l.ScrollHalfPageUp()
	case "<C-f>":
		l.ScrollPageDown()
	case "<C-b>":
		l.ScrollPageUp()
	case "g":
		if l.previousKey == "g" {
			l.ScrollTop()
		}
	case "<Home>":
		l.ScrollTop()
	case "G", "<End>":
		l.ScrollBottom()
	}

	if l.previousKey == "g" {
		l.previousKey = ""
	} else {
		l.previousKey = e.ID
	}
}

func (l *MyList) SetBorderStyle(style ui.Style) {
	l.BorderStyle = style
}

func (l *MyList) Refresh() {

}

func NewList(config map[string]interface{}, render chan<- struct{}) widget.Widget {
	l := widgets.NewList()
	l.Title = "List"
	l.Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[0] github.com/gizak/termui/v3",
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] baz",
	}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false

	result := &MyList{
		List: *l,
	}
	return result
}
