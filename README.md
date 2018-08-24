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

# create tada.json and copy below content into it
$ touch tada.json
```

### tada.json
```js
{
  "widgets": [{
      "name": "tada-github",
      "width": 3,
      "height": 3,
      "x": 1,
      "y":0
    }]
}
```

### Expected Widget
* GitHub pr (80% done!)
* Codeship status
* Circleci status
* Heroku status
* Metabase
* Rollbar

MIT