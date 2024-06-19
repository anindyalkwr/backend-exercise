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
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
)

type UserFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

func (f *UserFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "UserFetcher.Fetch")
	defer span.Finish()

	url := fmt.Sprintf("%suser?limit=10&page=%d", utils.GetBaseURL(), page)
	span.SetTag("url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		span.LogFields(log.Error(err))

		fmt.Println(err)
		return
	}

	req.Header.Set(appIDHeader, appId)

	resp, err := client.Do(req)
	if err != nil {
		span.LogFields(log.Error(err))

		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []models.User `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		span.LogFields(log.Error(err))

		fmt.Printf("Error parsing users response: %v\n", err)
		return
	}

	for _, user := range result.Data {
		userKey := fmt.Sprintf("user:%s", user.ID)
		userData := fmt.Sprintf("User name %s %s %s\n", user.Title, user.FirstName, user.LastName)

		fmt.Print(userData)

		err = f.Rdb.Set(ctx, userKey, userData, 0).Err()
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error storing user in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(userKey),
			Value: []byte(userData),
		})
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error sending user data to Kafka: %v\n", err)
		}
	}
}
