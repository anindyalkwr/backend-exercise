package fetcher

import (
	"backend-exercise/models"
	"backend-exercise/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
)

type PostFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

func (f *PostFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%spost?limit=10&page=%d", utils.GetBaseURL(), page)

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
		Data []models.Post `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error parsing posts response: %v\n", err)
		return
	}

	for _, post := range result.Data {
		fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.User.FirstName, post.User.LastName, post.Text, post.Likes, post.Tags, post.PublishDate)

		postData, err := json.Marshal(post)
		if err != nil {
			fmt.Printf("Error marshalling post data: %v\n", err)
			continue
		}

		postKey := fmt.Sprintf("post:%s", post.ID)
		err = f.Rdb.Set(context.Background(), postKey, postData, 0).Err()
		if err != nil {
			fmt.Printf("Error storing post in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Value: postData,
		})
		if err != nil {
			fmt.Printf("Error sending post data to Kafka: %v\n", err)
		}

	}
}
