package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetAppID() string {
	appId := os.Getenv("APP_ID")
	if appId == "" {
		log.Fatal("APP_ID not set in .env file")
	}
	return appId
}

func GetBaseURL() string {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		log.Fatal("BASE_URL not set in .env file")
	}
	return baseUrl
}

func GetKafkaURL() string {
	kafkaUrl := os.Getenv("KAFKA_URL")
	if kafkaUrl == "" {
		log.Fatal("KAFKA_URL not set in .env file")
	}
	return kafkaUrl
}

func GetKafkaTopic() string {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatal("KAFKA_TOPIC not set in .env file")
	}
	return kafkaTopic
}

func GetRedisURL() string {
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("REDIS_URL not set in .env file")
	}
	return redisUrl
}

func GetJaegerURL() string {
	jaegerUrl := os.Getenv("JAEGER_URL")
	if jaegerUrl == "" {
		log.Fatal("JAEGER_URL not set in .env file")
	}
	return jaegerUrl
}
