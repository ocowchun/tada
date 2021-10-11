package ghpr

import (
	"bytes"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	widget "github.com/ocowchun/tada/Widget"
	"github.com/shurcooL/githubv4"
	"log"
	"os/exec"
	"sync"
)

type ListType int8

const (
	PR ListType = iota
	Review
)

type GitHubPRList struct {
	previousKey string
	widgets.List
	pullRequests     []PullRequest
	prLock           sync.Mutex
	ghClient         *githubv4.Client
	renderCh         chan<- struct{}
	listType         ListType
	reviewerUsername string
}

func (l *GitHubPRList) HandleUIEvent(e ui.Event) {
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
		pr := l.selectedPullRequest()
		if pr != nil {
			cmd := exec.Command("open", pr.url)
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

func (l *GitHubPRList) SetBorderStyle(style ui.Style) {
	l.BorderStyle = style
}

type CIState int8

const (
	Pending CIState = iota
	Passed
	Failed
)

type PullRequest struct {
	title                 string
	ciState               CIState
	approvedCount         int32
	commentedCount        int32
	changesRequestedCount int32
	url                   string
}

func toRows(pullRequests []PullRequest) []string {
	result := make([]string, len(pullRequests))

	for i, pr := range pullRequests {
		var b bytes.Buffer
		if pr.ciState == Passed {
			b.WriteString("[V](fg:green) ")
		} else if pr.ciState == Failed {
			b.WriteString("[X](fg:red) ")
		} else {
			b.WriteString("[O](fg:yellow) ")
		}
		b.WriteString(pr.title)
		result[i] = b.String()
	}
	return result
}

func (l *GitHubPRList) fetchAndUpdatePullRequests() {
	result := make([]PullRequest, 0)
	var prs []GhPullRequest
	var err error

	if l.listType == Review {
		prs, err = fetchReviewRequests(l.ghClient, l.reviewerUsername)
	} else {
		prs, err = fetchPullRequests(l.ghClient)
	}
	if err != nil {
		log.Fatal(err)
	}

	for _, pr := range prs {
		ciState := Pending
		state := pr.Commits.Nodes[0].COMMIT.Status.State
		if state == githubv4.StatusStateSuccess {
			ciState = Passed
		} else if state == githubv4.StatusStateFailure {
			ciState = Failed
		}

		reviewStats := computePullRequestReviewStats(pr)
		result = append(result, PullRequest{
			title:                 pr.Title,
			url:                   pr.Url.String(),
			ciState:               ciState,
			approvedCount:         reviewStats[githubv4.PullRequestReviewStateApproved],
			commentedCount:        reviewStats[githubv4.PullRequestReviewStateCommented],
			changesRequestedCount: reviewStats[githubv4.PullRequestReviewStateChangesRequested],
		})
	}

	l.updatePullRequests(result)

}

func (l *GitHubPRList) Refresh() {
	l.updatePullRequests(make([]PullRequest, 0))
	l.Rows = []string{"loading..."}

	//log.Println("updatePR")
	go l.fetchAndUpdatePullRequests()
}

func (l *GitHubPRList) updatePullRequests(pullRequests []PullRequest) {
	l.prLock.Lock()
	defer l.prLock.Unlock()

	l.pullRequests = pullRequests
	l.Rows = toRows(pullRequests)
	l.renderCh <- struct{}{}
}

func (l *GitHubPRList) selectedPullRequest() *PullRequest {
	l.prLock.Lock()
	defer l.prLock.Unlock()

	i := l.SelectedRow
	if i >= len(l.pullRequests) {
		return nil
	}
	return &l.pullRequests[i]
}

func NewGitHubPRList(config map[string]interface{}, renderCh chan<- struct{}) widget.Widget {
	l := widgets.NewList()
	l.Title = "Pull Requests"
	l.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierUnderline)
	l.WrapText = false

	if _, existed := config["GITHUB_TOKEN"]; !existed {
		log.Fatal("You must provide GITHUB_TOKEN to enable this Widget")
	}

	accessToken := config["GITHUB_TOKEN"].(string)
	ghClient, err := newGitHubClient(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	listType := PR
	reviewerUsername := ""
	if str, existed := config["LIST_TYPE"]; existed {
		if str.(string) == "review" {
			listType = Review
			if config["GITHUB_USERNAME"] == nil {
				log.Fatal("You must provide GITHUB_USERNAME if LIST_TYPE = `review`")
			}
			l.Title = "Review Requests"
			reviewerUsername = config["GITHUB_USERNAME"].(string)
		}
	}

	component := &GitHubPRList{
		List:             *l,
		ghClient:         ghClient,
		renderCh:         renderCh,
		listType:         listType,
		reviewerUsername: reviewerUsername,
	}
	go component.Refresh()
	return component
}
