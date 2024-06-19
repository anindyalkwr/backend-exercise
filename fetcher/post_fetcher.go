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

type PostFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

func (f *PostFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "PostFetcher.Fetch")
	defer span.Finish()

	url := fmt.Sprintf("%spost?limit=10&page=%d", utils.GetBaseURL(), page)
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
		Data []models.Post `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		span.LogFields(log.Error(err))

		fmt.Printf("Error parsing posts response: %v\n", err)
		return
	}

	for _, post := range result.Data {
		postKey := fmt.Sprintf("post:%s", post.ID)
		postData := fmt.Sprintf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n\n", post.User.FirstName, post.User.LastName, post.Text, post.Likes, post.Tags, post.PublishDate)

		fmt.Print(postData)

		err = f.Rdb.Set(ctx, postKey, postData, 0).Err()
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error storing post in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(postKey),
			Value: []byte(postData),
		})
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error sending post data to Kafka: %v\n", err)
		}

	}
}
