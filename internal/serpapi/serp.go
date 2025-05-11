package serpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OrganicResult struct {
	Link string `json:"link"`
}

type SerpResponse struct {
	OrganicResults []OrganicResult `json:"organic_results"`
}

func GetTopUrls(apiKey string, query string, num int) ([]string, error) {
	url := fmt.Sprintf("https://serpapi.com/search.json?q=%s&engine=google&num=%d&api_key=%s", query, num, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SerpResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	var urls []string
	for _, item := range result.OrganicResults {
		urls = append(urls, item.Link)
	}

	return urls, nil
}
