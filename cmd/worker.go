package cmd

import (
	"backend-exercise/fetcher"
	"backend-exercise/utils"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/opentracing/opentracing-go"
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

	rdb := utils.InitRedisClient()

	kafkaWriter := utils.InitKafkaWriter()
	defer kafkaWriter.Close()

	tracer, closer := utils.InitJaegerTracer()
	if closer != nil {
		defer closer.Close()
	}
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
