package hackernews

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/ocowchun/tada/widget"
)

type HackerNewsBox struct {
	isHover        bool
	isFocus        bool
	width          int
	height         int
	stories        []*Story
	storyAmount    int
	loading        bool
	githubUsername string
	githubToken    string
	stopApp        func()
}

type Story struct {
	id           int `json:"id"`
	title        string
	isHover      bool
	commentCount int
}

func (s *Story) Url() string {
	return fmt.Sprintf("https://news.ycombinator.com/item?id=%v", s.id)
}

func (box *HackerNewsBox) Focus() {
	box.isFocus = true
}

func (box *HackerNewsBox) Blur() {
	box.isFocus = false
	storyIdx := findHoverStory(box.stories)
	if storyIdx != -1 {
		box.stories[storyIdx].isHover = false
	}
}

func (box *HackerNewsBox) Hover() {
	box.isHover = true
}

func (box *HackerNewsBox) Unhover() {
	box.isHover = false
}

func findHoverStory(prs []*Story) int {
	for i := 0; i < len(prs); i++ {
		if prs[i].isHover {
			return i
		}
	}
	return -1
}

func (box *HackerNewsBox) InputCaptureFactory(render func()) func(event *widget.KeyEvent) {
	return func(event *widget.KeyEvent) {
		switch event.Key {
		case widget.KeyRune:
			switch event.Rune {
			case 'r':
				box.loading = true
				render()
				box.stories = box.fetchStories()
				box.loading = false
				render()
			}
		case widget.KeyDown:
			storyIdx := findHoverStory(box.stories)
			if storyIdx == -1 {
				box.stories[0].isHover = true
			} else {
				box.stories[storyIdx].isHover = false
				newIdx := (storyIdx + 1) % len(box.stories)
				box.stories[newIdx].isHover = true
			}
			render()
		case widget.KeyUp:
			storyIdx := findHoverStory(box.stories)
			if storyIdx == -1 {
				box.stories[0].isHover = true
			} else {
				box.stories[storyIdx].isHover = false
				newIdx := (storyIdx - 1 + len(box.stories)) % len(box.stories)
				box.stories[newIdx].isHover = true
			}
			render()
		case widget.KeyEnter:
			storyIdx := findHoverStory(box.stories)
			if storyIdx != -1 {
				cmd := exec.Command("open", box.stories[storyIdx].Url())
				_, err := cmd.Output()
				if err != nil {
					log.Printf("Command finished with error: %v", err)
				}

			}
		}
	}
}

func (box *HackerNewsBox) Render(width int) []string {
	lines := []string{}
	if box.loading {
		line := &widget.Line{
			Width: width,
		}
		line.AddSentence(&widget.Sentence{Content: "loading...", Color: "white"})
		lines = append(lines, line.String())

	} else {
		for i := 0; i < len(box.stories); i++ {
			story := box.stories[i]
			line := &widget.Line{
				Width: width,
			}

			titleColor := "white"
			if story.isHover {
				titleColor = "red"
			}

			replacer := strings.NewReplacer("‘", "`",
				"’", "`",
				",", ",",
				"–", "-",
				"“", "\"",
				"”", "\"",
			)
			title := replacer.Replace(story.title)
			maxTitleLength := width - 10
			if maxTitleLength < 0 {
				maxTitleLength = 0
			}
			if len(title) > maxTitleLength {
				title = title[0:maxTitleLength]
			}

			line.AddSentence(&widget.Sentence{
				Content: title,
				Color:   titleColor,
			})

			lines = append(lines, line.String())
		}
	}
	return lines
}

func (box *HackerNewsBox) fetchStories() []*Story {
	hnStories, err := FetchStories(box.storyAmount)

	if err != nil {
		box.stopApp()
		fmt.Println("perform query failed:")
		fmt.Println(err.Error())
	}

	stories := []*Story{}
	for _, hnStory := range hnStories {
		story := &Story{
			id:           hnStory.ID,
			title:        hnStory.Title,
			commentCount: hnStory.CommentCount,
		}
		stories = append(stories, story)
	}
	return stories
}

func NewWidget(config widget.Config, stopApp func()) *widget.Widget {
	storyAmount := 8
	hnStoryAmount, ok := config.Options["HN_STORY_AMOUNT"].(int64)
	if ok {
		storyAmount = int(hnStoryAmount)
	}

	box := &HackerNewsBox{
		loading:     true,
		stopApp:     stopApp,
		storyAmount: storyAmount,
	}
	widget := widget.NewWidget(box)
	stories := []*Story{}
	box.stories = stories
	refreshInterval := 600
	go func() {
		for {
			box.stories = box.fetchStories()
			box.loading = false
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()
	return widget
}
