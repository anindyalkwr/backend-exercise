package fetcher

import (
	"backend-exercise/models"
	"backend-exercise/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type PostFetcher struct{}

func (f *PostFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%spost?limit=10&page=%d", utils.GetBaseURL(), page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set(appIDHeader, appId)
	q := req.URL.Query()
	q.Add("page", fmt.Sprintf("%d", page))
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []models.Post `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error parsing posts response: %v\n", err)
		return
	}

	for _, post := range result.Data {
		fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.User.FirstName, post.User.LastName, post.Text, post.Likes, post.Tags, post.PublishDate)
	}
}
