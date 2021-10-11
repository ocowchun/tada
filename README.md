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
name = "ghpr"
# widget width
width = 70
height = 8
# widget x-axis
x = 1
# widget y-axis
y = 1
  [widgets.Options]
    GITHUB_TOKEN = "your-secret-token"
    LIST_TYPE = "pr"
[[widgets]]
name = "ghpr"
# widget width
width = 70
height = 8
# widget x-axis
x = 1
# widget y-axis
y = 12
  [widgets.Options]
    GITHUB_TOKEN = "your-secret-token"
    LIST_TYPE = "review"
    GITHUB_USERNAME = "ocowchun"
```

### Expected Widget
* GitHub pr (80% done!)
* Circleci status
* Heroku status

MIT