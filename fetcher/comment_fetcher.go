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

type CommentFetcher struct {
	Rdb         *redis.Client
	KafkaWriter *kafka.Writer
}

func (f *CommentFetcher) Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup) {
	defer wg.Done()

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "CommentFetcher.Fetch")
	defer span.Finish()

	url := fmt.Sprintf("%scomment?limit=10&page=%d", utils.GetBaseURL(), page)
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
		Data []models.Comment `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		span.LogFields(log.Error(err))

		fmt.Printf("Error parsing comments response: %v\n", err)
		return
	}

	for _, comment := range result.Data {
		commentKey := fmt.Sprintf("comment:%s", comment.ID)
		commentData := fmt.Sprintf("Comment by %s %s:\n%s\n\nPost ID: %s\n", comment.User.FirstName, comment.User.LastName, comment.Message, comment.Post)

		fmt.Print(commentData)

		err = f.Rdb.Set(ctx, commentKey, commentData, 0).Err()
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error storing comment in Redis: %v\n", err)
			continue
		}

		err = f.KafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(commentKey),
			Value: []byte(commentData),
		})
		if err != nil {
			span.LogFields(log.Error(err))

			fmt.Printf("Error sending comment data to Kafka: %v\n", err)
		}
	}
}
