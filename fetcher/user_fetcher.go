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

type UserFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

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

		userData, err := json.Marshal(user)
		if err != nil {
			fmt.Printf("Error marshalling user data: %v\n", err)
			continue
		}

		userKey := fmt.Sprintf("user:%s", user.ID)
		err = f.Rdb.Set(context.Background(), userKey, userData, 0).Err()
		if err != nil {
			fmt.Printf("Error storing user in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Value: userData,
		})
		if err != nil {
			fmt.Printf("Error sending user data to Kafka: %v\n", err)
		}
	}
}
