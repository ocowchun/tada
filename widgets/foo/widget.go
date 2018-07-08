package foo

import (
	"github.com/gdamore/tcell"
	widget "github.com/ocowchun/tada/widget"
	"github.com/rivo/tview"
)

type FooWidget struct {
	isFocus bool
}

func (w *FooWidget) Focus(delegate func(p tview.Primitive)) {
	w.isFocus = true
}

func (w *FooWidget) Blur() {
	w.isFocus = false
}

func (w *FooWidget) InputCaptureFactory(render func()) func(event *tcell.EventKey) *tcell.EventKey {
	return nil
}

func (w *FooWidget) Render(width int) []string {
	strs := []string{}
	max := 3
	for i := 0; i < max; i++ {
		line := &widget.Line{
			Width: width,
		}
		line.AddSentence(&widget.Sentence{
			Content: "[]foo",
			Color:   "white",
		})
		strs = append(strs, line.String())
	}
	return strs
}

func NewWidget() *widget.Widget {
	box := &FooWidget{}
	widget := widget.NewWidget(box)
	return widget
}
