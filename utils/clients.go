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
		Brokers:  []string{GetKafkaURL()},
		Topic:    GetKafkaTopic(),
		Balancer: &kafka.Hash{},
	})
}

func InitKafkaReader(topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{GetKafkaURL()},
		Topic:   topic,
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
