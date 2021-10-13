package hackernews

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	widget "github.com/ocowchun/tada/Widget"
	"log"
	"os/exec"
	"sync"
)

type HackerNews struct {
	previousKey string
	widgets.List
	stories   []Story
	storyLock sync.Mutex
	renderCh  chan<- struct{}
}

func (l *HackerNews) HandleUIEvent(e ui.Event) {
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
	case "<Enter>":
		story := l.selectedStory()
		if story != nil {
			cmd := exec.Command("open", story.url)
			if _, err := cmd.Output(); err != nil {
				log.Println(err)
			}
		}
	}

	if l.previousKey == "g" {
		l.previousKey = ""
	} else {
		l.previousKey = e.ID
	}
}

func (l *HackerNews) SetBorderStyle(style ui.Style) {
	l.BorderStyle = style
}

type Story struct {
	title string
	url   string
}

func toRows(stories []Story) []string {
	result := make([]string, len(stories))

	for i, story := range stories {
		result[i] = fmt.Sprintf(" %v", story.title)
	}
	return result
}

func (l *HackerNews) fetchAndUpdateStories() {
	limit := 20
	stories, err := fetchStories(limit)
	if err != nil {
		l.Rows = []string{err.Error()}
		return
	}

	l.updateStories(stories)
}

func (l *HackerNews) Refresh() {
	l.updateStories(make([]Story, 0))
	l.Rows = []string{"loading..."}

	go l.fetchAndUpdateStories()
}

func (l *HackerNews) updateStories(stories []Story) {
	l.storyLock.Lock()
	defer l.storyLock.Unlock()

	l.stories = stories
	l.Rows = toRows(stories)
}

func (l *HackerNews) selectedStory() *Story {
	l.storyLock.Lock()
	defer l.storyLock.Unlock()

	i := l.SelectedRow
	if i >= len(l.stories) {
		return nil
	}
	return &l.stories[i]
}

func NewHackerNews(config map[string]interface{}, renderCh chan<- struct{}) widget.Widget {
	l := widgets.NewList()
	l.Title = "Hacker News"
	l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierUnderline)
	l.WrapText = false

	component := &HackerNews{
		List:     *l,
		renderCh: renderCh,
	}
	go component.Refresh()
	return component
}
