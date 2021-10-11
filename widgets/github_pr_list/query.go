package github_pr_list

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"os"
)

type GhReview struct {
	Author struct {
		Login string
	}
	State     githubv4.PullRequestReviewState
	CreatedAt githubv4.DateTime
}

type GhCommit struct {
	COMMIT struct {
		Status struct {
			State githubv4.StatusState
		}
	}
}

type ReviewRequestedEvent struct {
	CreatedAt         githubv4.DateTime
	RequestedReviewer struct {
		Typename string `graphql:"typename :__typename"`
		User     struct {
			Login string
		} `graphql:"... on User"`
	}
}

type GhTimelineItem struct {
	Typename string               `graphql:"typename :__typename"`
	Event    ReviewRequestedEvent `graphql:"... on ReviewRequestedEvent"`
}

type GhPullRequest struct {
	Title  string
	Url    githubv4.URI
	Author struct {
		Login string
	}
	Timeline struct {
		Nodes []GhTimelineItem
	} `graphql:"timeline(last:5)"`
	Repository struct {
		Name string
	}
	Commits struct {
		Nodes []GhCommit
	} `graphql:"commits(last:1)"`
	Reviews struct {
		Nodes []GhReview
	} `graphql:"reviews(last: 30)"`
}

func fetchPullRequests(client *githubv4.Client) ([]GhPullRequest, error) {
	var query struct {
		Viewer struct {
			Login        string
			PullRequests struct {
				Nodes []GhPullRequest
			} `graphql:"pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC})"`
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		return nil, err
	}
	return query.Viewer.PullRequests.Nodes, nil
}

func computePullRequestReviewStats(pr GhPullRequest) map[githubv4.PullRequestReviewState]int32 {
	reviews := pr.Reviews.Nodes
	m := make(map[string]GhReview)
	for _, review := range reviews {
		username := review.Author.Login
		if username == pr.Author.Login {
			continue
		}

		r, existed := m[username]
		if review.State == githubv4.PullRequestReviewStatePending || review.State == githubv4.PullRequestReviewStateDismissed {
			continue
		}

		if existed == false || r.CreatedAt.Before(review.CreatedAt.Time) {
			m[username] = review
		}
	}

	result := make(map[githubv4.PullRequestReviewState]int32)
	for _, review := range m {
		result[review.State]++
	}

	return result
}

func newGitHubClient() (*githubv4.Client, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	return client, nil
}
