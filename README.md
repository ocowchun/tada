## tada
Terminal dashboard 

Plugin architecture 
Golang

GitHub pr
Codeship status
Circleci status
Heroku status
Metabase
Rollbar 

## TODO
- [ ] package manager
- [ ] update packages

## package manager
```sh
tada install package-name
```

* read installed packages from a storage (file or kv store?)


### how to install packages?
install .so or install .go and compile them?

download source code from github and compile them

remove that plugin repo before install or maybe we can pull(fetch + merge) that repo

handle failed

## widget system
* select mode: focus between different widgets (i.e. unfocus widget A and focus widget B)
* function mode: perform specify widget function (i.e. scroll to next PR)
* make sure `keyboardIntercept` only work when view is focus
    * answer: yes
* don't extract widget interface untill write more than 2 widgets

## Github PR Plugins
Created
Review requests
focus on different grid
display 

I can press `enter` on each item and open it on browser


## render 
each line only care about content it doesn't care about border

widgest contains a lot line and border

implement first widget and then extract what can be extract
Implement first feature that is useful

https://github.com/shurcooL/githubv4
https://developer.github.com/v4/guides/forming-calls/#the-graphql-endpoint
https://developer.github.com/v4/object/pullrequestreview/
https://developer.github.com/v4/enum/pullrequestreviewstate/

```
{
  viewer {
    login
    name
    pullRequests(last: 10, states: [OPEN], orderBy: {field: CREATED_AT, direction: DESC}) {
      edges {
        node {
          title
          reviews(last: 5) {
            edges {
              node {
                author {
                  login
                }
                state
              }
            }
          }
        }
      }
    }
  }
}
```