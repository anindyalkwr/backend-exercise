package utils

import (
	"fmt"
	"io"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: GetRedisURL(),
	})
}

func InitKafkaWriter() *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "scraped-data",
		Balancer: &kafka.LeastBytes{},
	})
}

func InitKafkaReader() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "scraped-data",
	})
}

func InitJaegerTracer() (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		ServiceName: "backend-exercise",
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: GetJaegerURL(),
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		fmt.Printf("could not initialize jaeger tracer: %v\n", err)
		return nil, nil
	}
	return tracer, closer
}
