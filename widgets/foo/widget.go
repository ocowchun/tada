package foo

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/ocowchun/tada/widgets"
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
		line := &widgets.Line{
			Width: width,
		}
		line.AddSentence(&widgets.Sentence{
			Content: "[]foo",
			Color:   "white",
		})
		strs = append(strs, line.String())
	}
	return strs
}

func NewWidget() *widgets.Widget {
	box := &FooWidget{}
	widget := widgets.NewWidget(box)
	go func() {
		for {
			widget.Render()
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	return widget
}
