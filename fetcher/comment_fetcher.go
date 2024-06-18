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

type CommentFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

func (f *CommentFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	url := fmt.Sprintf("%scomment?limit=10&page=%d", utils.GetBaseURL(), page)

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
		Data []models.Comment `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error parsing comments response: %v\n", err)
		return
	}

	for _, comment := range result.Data {
		fmt.Printf("Comment by %s %s:\n%s\n\nPost ID: %s\n", comment.User.FirstName, comment.User.LastName, comment.Message, comment.Post)

		commentData, err := json.Marshal(comment)
		if err != nil {
			fmt.Printf("Error marshalling comment data: %v\n", err)
			continue
		}

		commentKey := fmt.Sprintf("comment:%s", comment.ID)
		err = f.Rdb.Set(context.Background(), commentKey, commentData, 0).Err()
		if err != nil {
			fmt.Printf("Error storing comment in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Value: commentData,
		})
		if err != nil {
			continue
			// fmt.Printf("Error sending comment data to Kafka: %v\n", err)
		}
	}
}
