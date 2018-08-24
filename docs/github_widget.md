# Github PR Plugins
Created
Review requests
focus on different grid
display

I can press `enter` on each item and open it on browser

## TODO
- [ ] Align PR no matter it has build status or not
- [ ] Diagnose abnormal exit status


## ref
https://github.com/shurcooL/githubv4
https://developer.github.com/v4/guides/forming-calls/#the-graphql-endpoint
https://developer.github.com/v4/object/pullrequestreview/
https://developer.github.com/v4/enum/pullrequestreviewstate/
using variables to decide PR count and review count
https://developer.github.com/v4/guides/forming-calls/#working-with-variables


```
{
 viewer {
    login
    name
    
    pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC}) {
      edges {
        node {
          title
          timeline(last:10){
            totalCount
            nodes{
              __typename
              ...on ReviewRequestedEvent{
                createdAt
                requestedReviewer{
                  ...on User{
                    login
                  }
                }
              }
              ...on Commit{
                message
              }
            }
          }
          reviewRequests(last:5){
            edges{
              node{
                requestedReviewer{
                  ...on User{
                    login
                  }
                }
              }
            }
          }
          commits(last:1){
            edges{
              node{
                commit {
                  commitUrl
                  status{
                    state
                  }
                }
              }
            }
          }
          reviews(last: 5,) {
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
        }
      }
    }
  }
}
```