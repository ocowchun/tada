package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	Url      string
	Response *http.Response
	Err      error
}

type HnStory struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	commentCount int    `json:"descendants"`
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

func fetchStoryIds(amount int) ([]int, error) {
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"
	res := sendRequest(url)

	if res.Err != nil {
		return nil, res.Err
	}
	respBody, _ := ioutil.ReadAll(res.Response.Body)
	var storyIds []int
	if err := json.Unmarshal(respBody, &storyIds); err != nil {
		return nil, err
	}
	return storyIds[0:amount], nil
}

func FetchStories(amount int) ([]HnStory, error) {
	storyIds, err := fetchStoryIds(amount)
	if err != nil {
		return nil, err
	}
	resCh := make(chan *HttpResponse)
	for _, storyID := range storyIds {
		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", storyID)
		go func() {
			resCh <- sendRequest(url)
		}()
	}

	storyMap := make(map[int]HnStory)

	for range storyIds {
		res := <-resCh
		if res.Err != nil {
			return nil, res.Err
		}
		respBody, err := ioutil.ReadAll(res.Response.Body)
		if err != nil {
			return nil, err
		}
		var story HnStory
		if err := json.Unmarshal(respBody, &story); err != nil {
			return nil, err
		}
		storyMap[story.ID] = story
	}
	stories := []HnStory{}
	for _, storyID := range storyIds {
		stories = append(stories, storyMap[storyID])
	}
	return stories, nil
}
