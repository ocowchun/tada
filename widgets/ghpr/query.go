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

type PullRequest struct {
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

func FetchPullRequestsWithGraphQL(client *ghbv4.Client) ([]*Issue, error) {
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
		return nil, err
	}
	issues := []*Issue{}
	for _, pr := range query.Viewer.PullRequests.Nodes {
		stateCountMap := util.ComputeReviewStatus(pr.Timeline.Nodes, pr.Reviews.Nodes, query.Viewer.Login)

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
	return issues, nil
}
