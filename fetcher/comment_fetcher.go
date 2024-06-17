package fetcher

import (
	"backend-exercise/models"
	"backend-exercise/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type CommentFetcher struct{}

func (f *CommentFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%scomment?limit=10&page=%d", utils.GetBaseURL(), page)

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
		Data []models.Comment `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error parsing comments response: %v\n", err)
		return
	}

	for _, comment := range result.Data {
		fmt.Printf("Comment by %s %s:\n%s\n\nPost ID: %s\n", comment.User.FirstName, comment.User.LastName, comment.Message, comment.Post)
	}
}
