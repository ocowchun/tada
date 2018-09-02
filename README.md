# SILLL IN DEVELOPMENT⛏

## tada
> Terminal dashboard

## Requeirements
* [Go](https://golang.org/)

## Install
```sh
$ go get -u github.com/ocowchun/tada

# Initialize tada required config
$ tada init

# Go to config folder and write config using your favorite editor
$ cd ~/.tada/
```

### tada.toml
```toml
[[widgets]]
# widget name
name = "tada-github-pr"
# widget width
width = 3
# widget height
height = 3
# widget x-axis
x = 1
# widget y-axis
y = 0
  [Widgets.Options]
    GITHUB_USERNAME = "your-github-username"
    GITHUB_TOKEN = "your-github-developer-token"
```

### Expected Widget
* GitHub pr (80% done!)
* Codeship status
* Circleci status
* Heroku status
* Metabase
* Rollbar

MIT