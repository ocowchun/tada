package ghpr

import (
	"context"

	"github.com/ocowchun/tada/widgets/util"
	ghbv4 "github.com/shurcooL/githubv4"
)

func FetchPullRequests(client *ghbv4.Client) ([]util.GhPullRequest, error) {
	var query struct {
		Viewer struct {
			Login        string
			PullRequests struct {
				Nodes []util.GhPullRequest
			} `graphql:"pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC})"`
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		return nil, err
	}
	return query.Viewer.PullRequests.Nodes, nil
}
