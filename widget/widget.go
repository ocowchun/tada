package widget

import (
	"math"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Box interface {
	Render(width int) []string
	InputCaptureFactory(render func()) func(event *KeyEvent)
	Focus()
	Blur()
}

func (w *Widget) Hover() {
	w.isHover = true
	w.Render()
}

func (w *Widget) Unhover() {
	w.isHover = false
	w.Render()
}

type Widget struct {
	box        Box
	Title      string
	textView   *tview.TextView
	isHover    bool
	isFocus    bool
	width      int
	height     int
	hasChanged bool
	sync.Mutex
	isDrawing bool
}

func (w *Widget) IsRendering() bool {
	return w.isDrawing
}

func NewWidget(box Box) *Widget {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)

	w := &Widget{
		box:      box,
		textView: textView,
		isHover:  false,
		isFocus:  false,
	}
	inputCapture := func(event *tcell.EventKey) *tcell.EventKey {
		keyEvent := ConvertEvent(event)
		box.InputCaptureFactory(w.Render)(keyEvent)
		return event
	}
	textView.SetInputCapture(inputCapture)
	return w
}

func (w *Widget) Draw(screen tcell.Screen) {
	w.Lock()
	if w.isDrawing {
		w.Unlock()
	} else {
		w.isDrawing = true
		w.Unlock()
		w.textView.Draw(screen)
		w.isDrawing = false
	}
}

func (w *Widget) GetRect() (int, int, int, int) {
	return w.textView.GetRect()
}

func (w *Widget) SetRect(x, y, width, height int) {
	w.textView.SetRect(x, y, width, height)
}

func (w *Widget) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return w.textView.InputHandler()
}

func (w *Widget) Focus(delegate func(p tview.Primitive)) {
	w.isFocus = true
	w.textView.Focus(delegate)
	if w.box.Focus != nil {
		w.box.Focus()
	}
	w.Render()
}

func (w *Widget) Blur() {
	w.isFocus = false
	w.textView.Blur()
	if w.box.Blur != nil {
		w.box.Blur()
	}
	w.Render()
}

func (w *Widget) GetFocusable() tview.Focusable {
	return w.textView.GetFocusable()
}

func buildSpace(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += " "
	}
	return line
}

func buildLine(text string, width int) string {
	// first char is ignore
	line := "" + text
	if len(line) > width {
		return line[0:width]
	} else {
		return line + buildSpace(width-len(line))
	}
}

func newTitleLine(title string, length int) string {
	titleLength := len(title)
	if titleLength >= length {
		return title[0:length]
	} else {
		l := float64(length - titleLength)
		leftLength := int(math.Ceil(l / float64(2)))
		rightLength := int(math.Floor(l / float64(2)))
		return newHorizontalLine(leftLength) + title + newHorizontalLine(rightLength)
	}
}

func newHorizontalLine(length int) string {
	line := ""
	max := length / 2
	for i := 0; i < max; i++ {
		line = line + "- "
	}
	if length%2 == 1 {
		line = line + "-"
	}
	return line
}

func (w *Widget) Render() {
	// stop set text when draw
	if !w.isDrawing {
		_, _, width, height := w.textView.GetRect()

		// I don't know why this broken
		// if w.width != width && w.height != height {
		w.width = width
		w.height = height
		// }
		lines := w.box.Render(w.width - 3)
		leftBorder := " |"
		rightBorder := "|"
		titleLine := " +" + newTitleLine(w.Title, width-3) + "+"
		horizontalLine := " +" + newHorizontalLine(width-3) + "+"

		if w.isFocus {
			titleLine = "[green]" + titleLine + "[white]"
			horizontalLine = "[green]" + horizontalLine + "[white]"
			leftBorder = "[green] |[white]"
			rightBorder = "[green]|[white]"
		} else if w.isHover {
			titleLine = "[yellow]" + titleLine + "[white]"
			horizontalLine = "[yellow]" + horizontalLine + "[white]"
			leftBorder = "[yellow] |[white]"
			rightBorder = "[yellow]|[white]"
		}

		if len(lines) > height-2 {
			lines = lines[0 : height-2]
		} else if len(lines) < height-2 {
			missingLineCount := height - (2 + len(lines))
			for i := 0; i < missingLineCount; i++ {
				lines = append(lines, buildLine("", width-3))
			}
		}
		text := titleLine

		for i := 0; i < len(lines); i++ {
			text += leftBorder + lines[i] + rightBorder
		}
		text = text + horizontalLine
		w.textView.SetText(text)
	}
}
