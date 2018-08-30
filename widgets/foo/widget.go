package foo

import (
	widget "github.com/ocowchun/tada/widget"
)

type FooWidget struct {
	isFocus bool
}

func (w *FooWidget) Focus() {
	w.isFocus = true
}

func (w *FooWidget) Blur() {
	w.isFocus = false
}

func (w *FooWidget) InputCaptureFactory(render func()) func(event *widget.KeyEvent) {
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

func NewWidget(config widget.Config) *widget.Widget {
	box := &FooWidget{}
	widget := widget.NewWidget(box)
	return widget
}
