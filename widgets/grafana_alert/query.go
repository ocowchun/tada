package grafana_alert

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

func sendRequest(apiKey string, url string) HttpResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", apiKey))
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

type AlertPayload struct {
	Name string `json:"name"`
	Path string `json:"url"`
}

func fetchAlerts(hostname string, apiKey string) ([]Alert, error) {
	alerts := make([]Alert, 0)

	url := fmt.Sprintf("%s/api/alerts?state=alerting", hostname)
	res := sendRequest(apiKey, url)

	if res.Err != nil {
		return alerts, res.Err
	}

	respBody, err := ioutil.ReadAll(res.Response.Body)
	if err != nil {
		return alerts, err
	}

	var alertPayloads []AlertPayload
	if err := json.Unmarshal(respBody, &alertPayloads); err != nil {
		return alerts, err
	}
	for _, pl := range alertPayloads {
		alerts = append(alerts, Alert{
			title: pl.Name,
			url:   fmt.Sprintf("%s%s", hostname, pl.Path),
		})
	}

	return alerts, nil
}
