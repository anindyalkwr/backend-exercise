package cmd

import (
	"backend-exercise/fetcher"
	"backend-exercise/utils"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run the worker to fetch data",
	Run: func(cmd *cobra.Command, args []string) {
		runWorker()
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}

func runWorker() {
	appId := utils.GetAppID()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "scraped-data",
	})
	defer kafkaWriter.Close()

	cfg := config.Configuration{
		ServiceName: "backend-exercise",
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Printf("could not initialize jaeger tracer: %v\n", err)
		return
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	var wg sync.WaitGroup
	client := &http.Client{}

	startTime := time.Now()

	fetchers := []fetcher.Fetcher{
		&fetcher.UserFetcher{Rdb: rdb, KafkaWriter: kafkaWriter},
		&fetcher.PostFetcher{Rdb: rdb, KafkaWriter: kafkaWriter},
		&fetcher.CommentFetcher{Rdb: rdb, KafkaWriter: kafkaWriter},
	}

	for i := 0; i < 10; i++ {
		for _, fetcher := range fetchers {
			wg.Add(1)
			go fetcher.Fetch(client, appId, i, &wg)
		}
	}

	wg.Wait()

	duration := time.Since(startTime)
	fmt.Printf("Total execution time: %v\n", duration)
}
