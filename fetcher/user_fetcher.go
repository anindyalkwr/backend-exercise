package fetcher

import (
	"backend-exercise/models"
	"backend-exercise/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UserFetcher struct{}

func (f *UserFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%suser?limit=10&page=%d", utils.GetBaseURL(), page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set(appIDHeader, appId)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []models.User `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error parsing users response: %v\n", err)
		return
	}

	for _, user := range result.Data {
		fmt.Printf("User name %s %s %s\n", user.Title, user.FirstName, user.LastName)
	}
}