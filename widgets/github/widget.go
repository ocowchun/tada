package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ocowchun/tada/widgets"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell"
	ghb "github.com/google/go-github/github"
	ghbv4 "github.com/shurcooL/githubv4"
)

type GitHubWidget struct {
	// box1    *tview.TextView
	isHover bool
	isFocus bool
	width   int
	height  int
	issues  []*Issue
}

type Issue struct {
	title                string
	isHover              bool
	url                  string
	status               string
	approvedCount        int
	changeRequestedCount int
	commentedCount       int
}

func (w *GitHubWidget) Focus(delegate func(p tview.Primitive)) {
	w.isFocus = true
}

func (w *GitHubWidget) Blur() {
	w.isFocus = false
	issueIdx := findHoverIssue(w.issues)
	if issueIdx != -1 {
		w.issues[issueIdx].isHover = false
	}
}

func (w *GitHubWidget) Hover() {
	w.isHover = true
}

func (w *GitHubWidget) Unhover() {
	w.isHover = false
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

// only care content no border
func (w *GitHubWidget) Render(width int) []string {
	lines := []string{}
	for i := 0; i < len(w.issues); i++ {
		issue := w.issues[i]
		line := &widgets.Line{
			Width: width,
		}

		if issue.isHover {
			line.AddSentence(&widgets.Sentence{Content: issue.title, Color: "red"})
		} else {
			line.AddSentence(&widgets.Sentence{Content: issue.title, Color: "white"})
		}

		if issue.approvedCount > 0 {
			line.AddSentence(&widgets.Sentence{
				Content: " V:" + strconv.Itoa(issue.approvedCount),
				Color:   "green",
			})
		}
		if issue.changeRequestedCount > 0 {
			line.AddSentence(&widgets.Sentence{
				Content: " X:" + strconv.Itoa(issue.changeRequestedCount),
				Color:   "red",
			})
		}
		if issue.commentedCount > 0 {
			line.AddSentence(&widgets.Sentence{
				Content: " C:" + strconv.Itoa(issue.commentedCount),
				Color:   "yellow",
			})
		}
		lines = append(lines, line.String())
	}

	return lines
}

func findHoverIssue(issues []*Issue) int {
	for i := 0; i < len(issues); i++ {
		if issues[i].isHover {
			return i
		}
	}
	return -1
}

func (w *GitHubWidget) InputCaptureFactory(render func()) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx == -1 {
				w.issues[0].isHover = true
			} else {
				w.issues[issueIdx].isHover = false
				newIdx := (issueIdx + 1) % len(w.issues)
				w.issues[newIdx].isHover = true
			}
			render()
		case tcell.KeyUp:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx == -1 {
				w.issues[0].isHover = true
			} else {
				w.issues[issueIdx].isHover = false
				newIdx := (issueIdx - 1 + len(w.issues)) % len(w.issues)
				w.issues[newIdx].isHover = true
			}
			render()
		case tcell.KeyEnter:
			issueIdx := findHoverIssue(w.issues)
			if issueIdx != -1 {
				cmd := exec.Command("open", w.issues[issueIdx].url)
				_, err := cmd.Output()
				if err != nil {
					log.Printf("Command finished with error: %v", err)
				}

			}
		}

		return event
	}
}

func initGithubClient() *ghb.Client {
	username := "ocowchun"
	password := os.Getenv("TADA_GITHUB_TOKEN")
	tp := &ghb.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	githubClient := ghb.NewClient(tp.Client())
	return githubClient
}

func initGithubV4Client() *ghbv4.Client {
	username := "ocowchun"
	password := os.Getenv("TADA_GITHUB_TOKEN")
	tp := &ghb.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client := ghbv4.NewClient(tp.Client())
	return client
}

func fetchPullRequestsWithGraphQL(client *ghbv4.Client) []*Issue {
	type review struct {
		Author struct {
			Login ghbv4.String
		}
		State ghbv4.PullRequestReviewState
	}
	type pullRequest struct {
		Title   string
		Url     ghbv4.URI
		Reviews struct {
			Nodes []review
		} `graphql:"reviews(last: 10)"`
	}

	var query struct {
		Viewer struct {
			Login        ghbv4.String
			Name         ghbv4.String
			CreatedAt    ghbv4.DateTime
			PullRequests struct {
				Nodes []pullRequest
			} `graphql:"pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC})"`
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
	}
	issues := []*Issue{}
	for _, pr := range query.Viewer.PullRequests.Nodes {
		stateCountMap := make(map[ghbv4.PullRequestReviewState]int)
		for _, r := range pr.Reviews.Nodes {
			stateCountMap[r.State] = stateCountMap[r.State] + 1
		}

		i := &Issue{
			title:                pr.Title,
			isHover:              false,
			url:                  pr.Url.String(),
			approvedCount:        stateCountMap[ghbv4.PullRequestReviewStateApproved],
			changeRequestedCount: stateCountMap[ghbv4.PullRequestReviewStateChangesRequested],
			commentedCount:       stateCountMap[ghbv4.PullRequestReviewStateCommented],
		}

		issues = append(issues, i)
	}
	return issues
}

func fetchPullRequests(client *ghb.Client) []*Issue {
	q := "author:ocowchun type:pr state:open"
	opts := &ghb.SearchOptions{}
	result, _, err := client.Search.Issues(context.Background(), q, opts)
	if err != nil {
		fmt.Println("Search.Repositories returned error: ", err)
	}
	issues := []*Issue{}
	for _, issue := range result.Issues {
		strs := strings.Split(issue.GetRepositoryURL(), "/")
		repoName := strs[len(strs)-1]
		i := &Issue{
			title:   repoName + "/" + issue.GetTitle(),
			isHover: false,
			url:     issue.GetHTMLURL(),
		}
		issues = append(issues, i)
	}
	return issues
}

func NewWidget() *widgets.Widget {
	box := &GitHubWidget{}
	widget := widgets.NewWidget(box)

	issues := []*Issue{}
	box.issues = issues
	refreshInterval := 20
	go func() {
		for {
			box.issues = fetchPullRequestsWithGraphQL(initGithubV4Client())
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()

	go func() {
		for {
			widget.Render()
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	return widget
}
