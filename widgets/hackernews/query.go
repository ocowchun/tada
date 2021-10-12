package hackernews

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type HttpResponse struct {
	Url      string
	Response *http.Response
	Err      error
}

func sendRequest(url string) HttpResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HttpResponse{
			Url: url,
			Err: err,
		}
	}

	res, err := client.Do(req)
	return HttpResponse{
		Url:      url,
		Response: res,
		Err:      err,
	}
}

func fetchStoryIds(limit int) ([]int, error) {
	var storyIds []int
	url := "https://hacker-news.firebaseio.com/v0/topstories.json"
	res := sendRequest(url)
	if res.Err != nil {
		return storyIds, res.Err
	}

	respBody, err := ioutil.ReadAll(res.Response.Body)
	if err != nil {
		return storyIds, err
	}

	if err := json.Unmarshal(respBody, &storyIds); err != nil {
		return storyIds, err
	}

	return storyIds[:limit], nil
}

type HnStory struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	CommentCount int    `json:"descendants"`
}

func fetchStories(limit int) ([]Story, error) {
	stories := make([]Story, 0)
	ids, err := fetchStoryIds(limit)
	if err != nil {
		return stories, err
	}

	concurrent := 3
	ch := make(chan HttpResponse, concurrent)
	var wg sync.WaitGroup
	wg.Add(len(ids))
	go func() {
		for _, id := range ids {
			go func(id int) {
				url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", id)
				ch <- sendRequest(url)
				wg.Done()
			}(id)
		}
	}()

	m := make(map[int]HnStory)

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		if res.Err != nil {
			continue
		}

		respBody, err := ioutil.ReadAll(res.Response.Body)
		if err != nil {
			continue
		}

		var story HnStory
		if err := json.Unmarshal(respBody, &story); err != nil {
			continue
		}
		m[story.ID] = story
	}

	for _, id := range ids {
		if h, existed := m[id]; existed {
			stories = append(stories, Story{
				title: h.Title,
				url:   fmt.Sprintf("https://news.ycombinator.com/item?id=%v", id),
			})
		}
	}

	return stories, nil
}
