package ghreview

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ocowchun/tada/widget"
	"github.com/ocowchun/tada/widgets/util"
	ghbv4 "github.com/shurcooL/githubv4"
)

type GhReviewBox struct {
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

func (box *GhReviewBox) Focus() {
	box.isFocus = true
}

func (box *GhReviewBox) Blur() {
	box.isFocus = false
	prIdx := findHoverPr(box.pullRequests)
	if prIdx != -1 {
		box.pullRequests[prIdx].isHover = false
	}
}

func (box *GhReviewBox) Hover() {
	box.isHover = true
}

func (box *GhReviewBox) Unhover() {
	box.isHover = false
}

func findHoverPr(prs []*PullRequest) int {
	for i := 0; i < len(prs); i++ {
		if prs[i].isHover {
			return i
		}
	}
	return -1
}

func (box *GhReviewBox) InputCaptureFactory(render func()) func(event *widget.KeyEvent) {
	return func(event *widget.KeyEvent) {
		switch event.Key {
		case widget.KeyRune:
			switch event.Rune {
			case 'r':
				box.loading = true
				render()
				box.pullRequests = box.fetchReviewRequests()
				box.loading = false
				render()
			}
		case widget.KeyDown:
			prIdx := findHoverPr(box.pullRequests)
			if len(box.pullRequests) > 0 {
				if prIdx == -1 {
					box.pullRequests[0].isHover = true
				} else {
					box.pullRequests[prIdx].isHover = false
					newIdx := (prIdx + 1) % len(box.pullRequests)
					box.pullRequests[newIdx].isHover = true
				}
				render()
			}
		case widget.KeyUp:
			if len(box.pullRequests) > 0 {
				prIdx := findHoverPr(box.pullRequests)
				if prIdx == -1 {
					box.pullRequests[0].isHover = true
				} else {
					box.pullRequests[prIdx].isHover = false
					newIdx := (prIdx - 1 + len(box.pullRequests)) % len(box.pullRequests)
					box.pullRequests[newIdx].isHover = true
				}
				render()
			}
		case widget.KeyEnter:
			if len(box.pullRequests) > 0 {
				prIdx := findHoverPr(box.pullRequests)
				if prIdx != -1 {
					cmd := exec.Command("open", box.pullRequests[prIdx].url)
					_, err := cmd.Output()
					if err != nil {
						log.Printf("Command finished with error: %v", err)
					}

				}
			}
		}
	}
}

func (box *GhReviewBox) Render(width int) []string {
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
			username := pr.authorUsername
			title := (pr.repositoryName + "/" + pr.title)
			maxTitleLength := width - 12 - len(username)
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
			line.AddSentence(&widget.Sentence{
				Content: " @" + username,
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

type PullRequest struct {
	title                string
	isHover              bool
	url                  string
	authorUsername       string
	status               ghbv4.StatusState
	approvedCount        int
	changeRequestedCount int
	commentedCount       int
	repositoryName       string
}

func (box *GhReviewBox) fetchReviewRequests() []*PullRequest {
	client := util.InitGithubV4Client(box.githubUsername, box.githubToken)
	ghPullRequests, err := FetchReviewRequests(client)
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
			authorUsername:       ghpr.Author.Login,
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
		fmt.Println(fmt.Sprintf("You must provide %v in tada.toml for tada-github-review", name))
		os.Exit(1)
	}
	return str
}

func NewWidget(config widget.Config, stopApp func()) *widget.Widget {
	githubUsername := getStringFromConfig(config, "GITHUB_USERNAME")
	githubToken := getStringFromConfig(config, "GITHUB_TOKEN")
	box := &GhReviewBox{
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
			box.pullRequests = box.fetchReviewRequests()
			box.loading = false
			widget.Render()
			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()
	return widget
}
