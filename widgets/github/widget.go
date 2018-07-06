package github

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
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
	isHover bool
	isFocus bool
	width   int
	height  int
	issues  []*Issue
	loading bool
}

type Issue struct {
	title                string
	isHover              bool
	url                  string
	status               ghbv4.StatusState
	approvedCount        int
	changeRequestedCount int
	commentedCount       int
	repositoryName       string
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

func (w *GitHubWidget) Render(width int) []string {
	lines := []string{}
	if w.loading {
		line := &widgets.Line{
			Width: width,
		}
		line.AddSentence(&widgets.Sentence{Content: "loading...", Color: "white"})
		lines = append(lines, line.String())

	} else {
		for i := 0; i < len(w.issues); i++ {
			issue := w.issues[i]
			line := &widgets.Line{
				Width: width,
			}
			switch issue.status {
			case ghbv4.StatusStateSuccess:
				line.AddSentence(&widgets.Sentence{Content: "V ", Color: "green"})
			case ghbv4.StatusStatePending:
				line.AddSentence(&widgets.Sentence{Content: "O ", Color: "yellow"})
			case ghbv4.StatusStateFailure:
				line.AddSentence(&widgets.Sentence{Content: "X ", Color: "red"})

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

			line.AddSentence(&widgets.Sentence{
				Content: title,
				Color:   titleColor,
			})

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
		case tcell.KeyRune:
			switch event.Rune() {
			case 'r':
				w.loading = true
				render()
				w.issues = fetchPullRequestsWithGraphQL(initGithubV4Client())
				w.loading = false
				render()
			}
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

func initGithubV4Client() *ghbv4.Client {
	username := os.Getenv("TADA_GITHUB_USERNAME")
	password := os.Getenv("TADA_GITHUB_TOKEN")
	tp := &ghb.BasicAuthTransport{
		Username: username,
		Password: password,
	}
	client := ghbv4.NewClient(tp.Client())
	return client
}

type Review struct {
	Author struct {
		Login string
	}
	State     ghbv4.PullRequestReviewState
	CreatedAt ghbv4.DateTime
}

type Commit struct {
	COMMIT struct {
		Status struct {
			State ghbv4.StatusState
		}
	}
}

type ReviewRequestedEvent struct {
	CreatedAt         ghbv4.DateTime
	RequestedReviewer struct {
		Typename string `graphql:"typename :__typename"`
		User     struct {
			Login string
		} `graphql:"... on User"`
	}
}
type TimelineItem struct {
	Typename string               `graphql:"typename :__typename"`
	Event    ReviewRequestedEvent `graphql:"... on ReviewRequestedEvent"`
}

type PullRequest struct {
	Title    string
	Url      ghbv4.URI
	Timeline struct {
		Nodes []TimelineItem
	} `graphql:"timeline(last:5)"`
	Repository struct {
		Name string
	}
	Commits struct {
		Nodes []Commit
	} `graphql:"commits(last:1)"`
	Reviews struct {
		Nodes []Review
	} `graphql:"reviews(last: 10)"`
}

func fetchPullRequestsWithGraphQL(client *ghbv4.Client) []*Issue {

	var query struct {
		Viewer struct {
			Login        string
			Name         ghbv4.String
			CreatedAt    ghbv4.DateTime
			PullRequests struct {
				Nodes []PullRequest
			} `graphql:"pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC})"`
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		fmt.Println(err)
	}
	issues := []*Issue{}
	for _, pr := range query.Viewer.PullRequests.Nodes {
		stateCountMap := computeReviewStatus(pr, query.Viewer.Login)
		i := &Issue{
			title:                pr.Title,
			isHover:              false,
			url:                  pr.Url.String(),
			approvedCount:        stateCountMap[ghbv4.PullRequestReviewStateApproved],
			changeRequestedCount: stateCountMap[ghbv4.PullRequestReviewStateChangesRequested],
			commentedCount:       stateCountMap[ghbv4.PullRequestReviewStateCommented],
			status:               pr.Commits.Nodes[0].COMMIT.Status.State,
			repositoryName:       pr.Repository.Name,
		}

		issues = append(issues, i)
	}
	return issues
}

type ReviewEvent struct {
	username  string
	createdAt ghbv4.DateTime
	action    string
}

type ByCreatedAt []ReviewEvent

func (a ByCreatedAt) Len() int {
	return len(a)
}

func (a ByCreatedAt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByCreatedAt) Less(i, j int) bool { return a[i].createdAt.Unix() < a[j].createdAt.Unix() }

func computeReviewStatus(pr PullRequest, authorUsername string) map[ghbv4.PullRequestReviewState]int {
	stateCountMap := make(map[ghbv4.PullRequestReviewState]int)
	complexWay := true
	if complexWay {
		reviewEvents := []ReviewEvent{}
		//idea: request:1,request:2, change:1,request:1,approve:1 =>approve:1, request:2
		for _, event := range pr.Timeline.Nodes {
			if event.Typename == "ReviewRequestedEvent" {
				evt := ReviewEvent{
					username:  event.Event.RequestedReviewer.User.Login,
					createdAt: event.Event.CreatedAt,
					action:    "requested",
				}
				reviewEvents = append(reviewEvents, evt)
			}
		}

		for _, review := range pr.Reviews.Nodes {
			if review.Author.Login != authorUsername {
				action := ""
				switch review.State {
				case ghbv4.PullRequestReviewStateApproved:
					action = "approved"
				case ghbv4.PullRequestReviewStateChangesRequested:
					action = "changesRequested"
				case ghbv4.PullRequestReviewStateCommented:
					action = "commented"
				}

				evt := ReviewEvent{
					username:  review.Author.Login,
					createdAt: review.CreatedAt,
					action:    action,
				}
				reviewEvents = append(reviewEvents, evt)
			}
		}
		sort.Sort(ByCreatedAt(reviewEvents))
		reviewMap := make(map[string]string)
		for _, event := range reviewEvents {
			reviewMap[event.username] = event.action
		}
		stateApproved := ghbv4.PullRequestReviewStateApproved
		stateChangesRequested := ghbv4.PullRequestReviewStateChangesRequested
		stateCommented := ghbv4.PullRequestReviewStateCommented
		for _, action := range reviewMap {
			switch action {
			case "approved":
				stateCountMap[stateApproved] = stateCountMap[stateApproved] + 1
			case "changesRequested":
				stateCountMap[stateChangesRequested] = stateCountMap[stateChangesRequested] + 1
			case "commented":
				stateCountMap[stateCommented] = stateCountMap[stateCommented] + 1
			}
		}

	} else {
		reviewrMap := make(map[string]bool)
		for _, r := range pr.Reviews.Nodes {
			if reviewrMap[r.Author.Login] == false {
				reviewrMap[r.Author.Login] = true
				stateCountMap[r.State] = stateCountMap[r.State] + 1
			}
		}
	}

	return stateCountMap
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
	box := &GitHubWidget{loading: true}
	widget := widgets.NewWidget(box)

	issues := []*Issue{}
	box.issues = issues
	refreshInterval := 120
	go func() {
		for {
			box.issues = fetchPullRequestsWithGraphQL(initGithubV4Client())
			box.loading = false
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()
	return widget
}
