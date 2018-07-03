package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Sentence struct {
	Content string
	Color   string
}

func (s *Sentence) len() int {
	return len(s.Content)
}

func (s *Sentence) render(maxLength int) string {
	text := s.Content
	if s.len() > maxLength {
		text = s.Content[0:maxLength]
	}
	if s.Color != "" && s.Color != "white" {
		text = "[" + s.Color + "]" + text + "[white]"
	}
	return text
}

type Line struct {
	sentences []*Sentence
	Width     int
}

func (l *Line) AddSentence(s *Sentence) {
	l.sentences = append(l.sentences, s)
}

func (l *Line) String() string {
	maxWidth := l.Width
	currentWidth := 0
	result := ""
	for _, sentence := range l.sentences {
		if currentWidth+sentence.len() <= maxWidth {
			currentWidth = currentWidth + sentence.len()
			result += sentence.render(sentence.len())
		}
	}
	return result + buildSpace(maxWidth-currentWidth)
}

type Box interface {
	Render(width int) []string
	InputCaptureFactory(render func()) func(event *tcell.EventKey) *tcell.EventKey
	Focus(delegate func(p tview.Primitive))
	Blur()
	// Hover()
	// Unhover()
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
	box      Box
	textView *tview.TextView
	isHover  bool
	isFocus  bool
	width    int
	height   int
}

func NewWidget(box Box) *Widget {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	// box1.SetBorder(false)

	w := &Widget{
		box:      box,
		textView: textView,
		isHover:  false,
		isFocus:  false,
	}
	textView.SetInputCapture(box.InputCaptureFactory(w.Render))
	return w
}

func (w *Widget) Draw(screen tcell.Screen) {
	w.textView.Draw(screen)
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
		w.box.Focus(delegate)
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
	_, _, width, height := w.textView.GetRect()

	// I don't know why this broken
	// if w.width != width && w.height != height {
	w.width = width
	w.height = height
	// }
	lines := w.box.Render(w.width - 3)
	leftBorder := " |"
	rightBorder := "|"
	horizontalLine := " +" + newHorizontalLine(width-3) + "+"

	if w.isFocus {
		horizontalLine = "[green] +" + newHorizontalLine(width-3) + "+[white]"
		leftBorder = "[green] |[white]"
		rightBorder = "[green]|[white]"
	} else if w.isHover {
		horizontalLine = "[yellow] +" + newHorizontalLine(width-3) + "+[white]"
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
	text := horizontalLine

	for i := 0; i < len(lines); i++ {
		text += leftBorder + lines[i] + rightBorder
	}
	text = text + horizontalLine
	w.textView.SetText(text)
}
