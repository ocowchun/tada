package ghpr

import (
	"context"

	util "github.com/ocowchun/tada/widgets/util"
	ghbv4 "github.com/shurcooL/githubv4"
)

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

type GhPullRequest struct {
	Title    string
	Url      ghbv4.URI
	Timeline struct {
		Nodes []util.GhTimelineItem
	} `graphql:"timeline(last:5)"`
	Repository struct {
		Name string
	}
	Commits struct {
		Nodes []util.GhCommit
	} `graphql:"commits(last:1)"`
	Reviews struct {
		Nodes []util.GhReview
	} `graphql:"reviews(last: 10)"`
}

func FetchPullRequestsWithGraphQL(client *ghbv4.Client) ([]*PullRequest, error) {
	var query struct {
		Viewer struct {
			Login        string
			Name         ghbv4.String
			CreatedAt    ghbv4.DateTime
			PullRequests struct {
				Nodes []GhPullRequest
			} `graphql:"pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC})"`
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		return nil, err
	}
	pullRequests := []*PullRequest{}
	for _, ghpr := range query.Viewer.PullRequests.Nodes {
		stateCountMap := util.ComputeReviewStatus(ghpr.Timeline.Nodes, ghpr.Reviews.Nodes, query.Viewer.Login)

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
	return pullRequests, nil
}
