package widget

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
