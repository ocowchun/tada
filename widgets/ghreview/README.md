V my-repo/the-branch-name @author_username V:2 

## GraphQL query

```
query($review_query: String!)
{
  viewer {
    login
  }
  search(last: 10, type: ISSUE, query: $review_query) {
    edges {
      node {
        ... on PullRequest {
          title
          author {
            login
          }
          reviews(last: 20) {
            edges {
              node {
                author {
                  login
                }
                state
                createdAt
              }
            }
          }
          commits(last: 1) {
            edges {
              node {
                commit {
                  status {
                    state
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}

```

```
{
    "review_query": "is:open is:pr review-requested:ocowchun archived:false"
}
```