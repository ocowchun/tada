package util

import (
	"sort"

	ghbv4 "github.com/shurcooL/githubv4"
)

type GhReview struct {
	Author struct {
		Login string
	}
	State     ghbv4.PullRequestReviewState
	CreatedAt ghbv4.DateTime
}

type GhCommit struct {
	COMMIT struct {
		Status struct {
			State ghbv4.StatusState
		}
	}
}

type ReviewEvent struct {
	username  string
	createdAt ghbv4.DateTime
	action    string
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
type GhTimelineItem struct {
	Typename string               `graphql:"typename :__typename"`
	Event    ReviewRequestedEvent `graphql:"... on ReviewRequestedEvent"`
}

// type PullRequest struct {
// 	Title    string
// 	Url      ghbv4.URI
// 	Timeline struct {
// 		Nodes []TimelineItem
// 	} `graphql:"timeline(last:5)"`
// 	Repository struct {
// 		Name string
// 	}
// 	Commits struct {
// 		Nodes []GhCommit
// 	} `graphql:"commits(last:1)"`
// 	Reviews struct {
// 		Nodes []GhReview
// 	} `graphql:"reviews(last: 10)"`
// }

type ByCreatedAt []ReviewEvent

func (a ByCreatedAt) Len() int {
	return len(a)
}

func (a ByCreatedAt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByCreatedAt) Less(i, j int) bool { return a[i].createdAt.Unix() < a[j].createdAt.Unix() }

func ComputeReviewStatus(timelineItems []GhTimelineItem, reviews []GhReview, authorUsername string) map[ghbv4.PullRequestReviewState]int {
	stateCountMap := make(map[ghbv4.PullRequestReviewState]int)
	reviewEvents := []ReviewEvent{}
	//idea: request:1,request:2, change:1,request:1,approve:1 =>approve:1, request:2
	for _, event := range timelineItems {
		if event.Typename == "ReviewRequestedEvent" {
			evt := ReviewEvent{
				username:  event.Event.RequestedReviewer.User.Login,
				createdAt: event.Event.CreatedAt,
				action:    "requested",
			}
			reviewEvents = append(reviewEvents, evt)
		}
	}

	for _, review := range reviews {
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

	return stateCountMap
}
