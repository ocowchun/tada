package ghpr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	widget "github.com/ocowchun/tada/widget"

	ghb "github.com/google/go-github/github"
	ghbv4 "github.com/shurcooL/githubv4"
)

type GitHubBox struct {
	isHover        bool
	isFocus        bool
	width          int
	height         int
	issues         []*PullRequest
	loading        bool
	githubUsername string
	githubToken    string
	stopApp        func()
}

type PullRequest struct {
	title                string
	isHover              bool
	url                  string
	status               ghbv4.StatusState
	approvedCount        int
	changeRequestedCount int
	commentedCount       int
	repositoryName       string
}

func (w *GitHubBox) Focus() {
	w.isFocus = true
}

func (w *GitHubBox) Blur() {
	w.isFocus = false
	issueIdx := findHoverIssue(w.issues)
	if issueIdx != -1 {
		w.issues[issueIdx].isHover = false
	}
}

func (w *GitHubBox) Hover() {
	w.isHover = true
}

func (w *GitHubBox) Unhover() {
	w.isHover = false
}

func (w *GitHubBox) Render(width int) []string {
	lines := []string{}
	if w.loading {
		line := &widget.Line{
			Width: width,
		}
		line.AddSentence(&widget.Sentence{Content: "loading...", Color: "white"})
		lines = append(lines, line.String())

	} else {
		for i := 0; i < len(w.issues); i++ {
			issue := w.issues[i]
			line := &widget.Line{
				Width: width,
			}

			switch issue.status {
			case ghbv4.StatusStateSuccess:
				line.AddSentence(&widget.Sentence{Content: "V ", Color: "green"})
			case ghbv4.StatusStatePending:
				line.AddSentence(&widget.Sentence{Content: "O ", Color: "yellow"})
			case ghbv4.StatusStateFailure:
				line.AddSentence(&widget.Sentence{Content: "X ", Color: "red"})
			case "":
				line.AddSentence(&widget.Sentence{Content: "  ", Color: "white"})
			}

			titleColor := "white"
			if issue.isHover {
				titleColor = "red"
			}

			title := (issue.repositoryName + "/" + issue.title)
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

			if issue.approvedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " V:" + strconv.Itoa(issue.approvedCount),
					Color:   "green",
				})
			}
			if issue.changeRequestedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " X:" + strconv.Itoa(issue.changeRequestedCount),
					Color:   "red",
				})
			}
			if issue.commentedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " C:" + strconv.Itoa(issue.commentedCount),
					Color:   "yellow",
				})
			}
			lines = append(lines, line.String())
		}
	}

	return lines
}

func findHoverIssue(issues []*PullRequest) int {
	for i := 0; i < len(issues); i++ {
		if issues[i].isHover {
			return i
		}
	}
	return -1
}

func (w *GitHubBox) InputCaptureFactory(render func()) func(event *widget.KeyEvent) {
	return func(event *widget.KeyEvent) {
		switch event.Key {
		case widget.KeyRune:
			switch event.Rune {
			case 'r':
				w.loading = true
				render()
				w.issues = w.fetchPullRequestsWithGraphQL(w.initGithubV4Client())
				w.loading = false
				render()
			}
		case widget.KeyDown:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx == -1 {
				w.issues[0].isHover = true
			} else {
				w.issues[issueIdx].isHover = false
				newIdx := (issueIdx + 1) % len(w.issues)
				w.issues[newIdx].isHover = true
			}
			render()
		case widget.KeyUp:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx == -1 {
				w.issues[0].isHover = true
			} else {
				w.issues[issueIdx].isHover = false
				newIdx := (issueIdx - 1 + len(w.issues)) % len(w.issues)
				w.issues[newIdx].isHover = true
			}
			render()
		case widget.KeyEnter:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx != -1 {
				cmd := exec.Command("open", w.issues[issueIdx].url)
				_, err := cmd.Output()
				if err != nil {
					log.Printf("Command finished with error: %v", err)
				}

			}
		}
	}
}

func (w *GitHubBox) initGithubV4Client() *ghbv4.Client {
	tp := &ghb.BasicAuthTransport{
		Username: w.githubUsername,
		Password: w.githubToken,
	}
	client := ghbv4.NewClient(tp.Client())
	return client
}

func (w *GitHubBox) fetchPullRequestsWithGraphQL(client *ghbv4.Client) []*PullRequest {
	pullRequests, err := FetchPullRequestsWithGraphQL(client)
	if err != nil {
		w.stopApp()
		fmt.Println("perform query failed:")
		fmt.Println(err.Error())
	}
	return pullRequests
}

func getStringFromConfig(config widget.Config, name string) string {
	str, ok := config.Options[name].(string)
	if !ok {
		fmt.Println(fmt.Sprintf("You must provide %v in tada.toml for tada-github", name))
		os.Exit(1)
	}
	return str
}

func NewWidget(config widget.Config, stopApp func()) *widget.Widget {
	githubUsername := getStringFromConfig(config, "GITHUB_USERNAME")
	githubToken := getStringFromConfig(config, "GITHUB_TOKEN")
	box := &GitHubBox{
		loading:        true,
		githubUsername: githubUsername,
		githubToken:    githubToken,
		stopApp:        stopApp,
	}
	widget := widget.NewWidget(box)

	issues := []*PullRequest{}
	box.issues = issues
	refreshInterval := 120
	go func() {
		for {
			box.issues = box.fetchPullRequestsWithGraphQL(box.initGithubV4Client())
			box.loading = false
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()
	return widget
}
