package foo

import (
	widget "github.com/ocowchun/tada/widget"
)

type FooBox struct {
	isFocus bool
}

func (w *FooBox) Focus() {
	w.isFocus = true
}

func (w *FooBox) Blur() {
	w.isFocus = false
}

func (w *FooBox) InputCaptureFactory(render func()) func(event *widget.KeyEvent) {
	return nil
}

func (w *FooBox) Render(width int) []string {
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

func NewWidget(config widget.Config, stopApp func()) *widget.Widget {
	box := &FooBox{}
	widget := widget.NewWidget(box)
	return widget
}
