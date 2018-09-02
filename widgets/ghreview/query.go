package ghreview

import (
	"context"

	"github.com/ocowchun/tada/widgets/util"
	ghbv4 "github.com/shurcooL/githubv4"
)

type PullRequestItem struct {
	Typename string             `graphql:"typename :__typename"`
	Pr       util.GhPullRequest `graphql:"... on PullRequest"`
}

func FetchReviewRequests(client *ghbv4.Client) ([]util.GhPullRequest, error) {
	githubUsername := "ocowchun"
	variables := map[string]interface{}{
		"review_query": ghbv4.String("is:open is:pr review-requested:" + githubUsername + " archived:false"),
	}

	var query struct {
		Search struct {
			Nodes []PullRequestItem
		} `graphql:"search(last: 10, type: ISSUE, query: $review_query)"`
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}
	pullRequests := []util.GhPullRequest{}

	for _, node := range query.Search.Nodes {
		pullRequests = append(pullRequests, node.Pr)
	}
	return pullRequests, nil
}
