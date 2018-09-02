package ghpr

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	widget "github.com/ocowchun/tada/widget"
	"github.com/ocowchun/tada/widgets/util"

	ghbv4 "github.com/shurcooL/githubv4"
)

type GitHubBox struct {
	isHover        bool
	isFocus        bool
	width          int
	height         int
	pullRequests   []*PullRequest
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

func (box *GitHubBox) Focus() {
	box.isFocus = true
}

func (box *GitHubBox) Blur() {
	box.isFocus = false
	prIdx := findHoverPr(box.pullRequests)
	if prIdx != -1 {
		box.pullRequests[prIdx].isHover = false
	}
}

func (box *GitHubBox) Hover() {
	box.isHover = true
}

func (box *GitHubBox) Unhover() {
	box.isHover = false
}

func (box *GitHubBox) Render(width int) []string {
	lines := []string{}
	if box.loading {
		line := &widget.Line{
			Width: width,
		}
		line.AddSentence(&widget.Sentence{Content: "loading...", Color: "white"})
		lines = append(lines, line.String())

	} else {
		for i := 0; i < len(box.pullRequests); i++ {
			pr := box.pullRequests[i]
			line := &widget.Line{
				Width: width,
			}

			switch pr.status {
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
			if pr.isHover {
				titleColor = "red"
			}

			title := (pr.repositoryName + "/" + pr.title)
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

			if pr.approvedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " V:" + strconv.Itoa(pr.approvedCount),
					Color:   "green",
				})
			}
			if pr.changeRequestedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " X:" + strconv.Itoa(pr.changeRequestedCount),
					Color:   "red",
				})
			}
			if pr.commentedCount > 0 {
				line.AddSentence(&widget.Sentence{
					Content: " C:" + strconv.Itoa(pr.commentedCount),
					Color:   "yellow",
				})
			}
			lines = append(lines, line.String())
		}
	}

	return lines
}

func findHoverPr(pullRequests []*PullRequest) int {
	for i := 0; i < len(pullRequests); i++ {
		if pullRequests[i].isHover {
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
				w.pullRequests = w.fetchPullRequestsWithGraphQL()
				w.loading = false
				render()
			}
		case widget.KeyDown:
			prIdx := findHoverPr(w.pullRequests)
			if prIdx == -1 {
				w.pullRequests[0].isHover = true
			} else {
				w.pullRequests[prIdx].isHover = false
				newIdx := (prIdx + 1) % len(w.pullRequests)
				w.pullRequests[newIdx].isHover = true
			}
			render()
		case widget.KeyUp:
			prIdx := findHoverPr(w.pullRequests)
			if prIdx == -1 {
				w.pullRequests[0].isHover = true
			} else {
				w.pullRequests[prIdx].isHover = false
				newIdx := (prIdx - 1 + len(w.pullRequests)) % len(w.pullRequests)
				w.pullRequests[newIdx].isHover = true
			}
			render()
		case widget.KeyEnter:
			prIdx := findHoverPr(w.pullRequests)
			if prIdx != -1 {
				cmd := exec.Command("open", w.pullRequests[prIdx].url)
				_, err := cmd.Output()
				if err != nil {
					log.Printf("Command finished with error: %v", err)
				}

			}
		}
	}
}

func (box *GitHubBox) fetchPullRequestsWithGraphQL() []*PullRequest {
	client := util.InitGithubV4Client(box.githubUsername, box.githubToken)
	ghPullRequests, err := FetchPullRequests(client)
	if err != nil {
		box.stopApp()
		fmt.Println("perform query failed:")
		fmt.Println(err.Error())
	}

	pullRequests := []*PullRequest{}
	for _, ghpr := range ghPullRequests {
		stateCountMap := util.ComputeReviewStatus(ghpr.Timeline.Nodes,
			ghpr.Reviews.Nodes, ghpr.Author.Login)
		pr := &PullRequest{
			title:                ghpr.Title,
			isHover:              false,
			url:                  ghpr.Url.String(),
			approvedCount:        stateCountMap[ghbv4.PullRequestReviewStateApproved],
			changeRequestedCount: stateCountMap[ghbv4.PullRequestReviewStateChangesRequested],
			commentedCount:       stateCountMap[ghbv4.PullRequestReviewStateCommented],
			status:               ghpr.Commits.Nodes[0].COMMIT.Status.State,
			repositoryName:       ghpr.Repository.Name,
		}

		pullRequests = append(pullRequests, pr)
	}
	return pullRequests
}

func getStringFromConfig(config widget.Config, name string) string {
	str, ok := config.Options[name].(string)
	if !ok {
		fmt.Println(fmt.Sprintf("You must provide %v in tada.toml for tada-github-pr", name))
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

	pullRequests := []*PullRequest{}
	box.pullRequests = pullRequests
	refreshInterval := 120
	go func() {
		for {
			box.pullRequests = box.fetchPullRequestsWithGraphQL()
			box.loading = false
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()
	return widget
}
