package hackernews

import (
	"net/http"
)

type HttpResponse struct {
	Url      string
	Response *http.Response
	Err      error
}

type Story struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
}

func sendRequest(url string) *HttpResponse {
	// Request (GET https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty)

	// Create client
	client := &http.Client{}
	// "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"
	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &HttpResponse{url, nil, err}
	}

	// Fetch Request
	resp, err := client.Do(req)
	return &HttpResponse{url, resp, err}
}

func FetchStories() {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"

}
