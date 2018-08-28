# SILLL IN DEVELOPMENT⛏

## tada
> Terminal dashboard

## Requeirements
* [Go](https://golang.org/)

## Install
```sh
$ go get -u github.com/ocowchun/tada
$ export TADA_GITHUB_TOKEN="your-personal-token"
$ mkdir ~/.tada

# create tada.toml and copy below content into it
$ touch tada.toml
```

### tada.toml
```toml
[[widgets]]
# widget name
name = "tada-github"
# widget width
width = 3
# widget height
height = 3
# widget x-axis
x = 1
# widget y-axis
y = 0
```

### Expected Widget
* GitHub pr (80% done!)
* Codeship status
* Circleci status
* Heroku status
* Metabase
* Rollbar

MIT