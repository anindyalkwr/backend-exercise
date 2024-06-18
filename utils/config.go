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
	appId := os.Getenv("BASE_URL")
	if appId == "" {
		log.Fatal("BASE_URL not set in .env file")
	}
	return appId
}

func GetKafkaURL() string {
	appId := os.Getenv("KAFKA_URL")
	if appId == "" {
		log.Fatal("KAFKA_URL not set in .env file")
	}
	return appId
}

func GetRedisURL() string {
	appId := os.Getenv("REDIS_URL")
	if appId == "" {
		log.Fatal("REDIS_URL not set in .env file")
	}
	return appId
}

func GetJaegerURL() string {
	appId := os.Getenv("JAEGER_URL")
	if appId == "" {
		log.Fatal("JAEGER_URL not set in .env file")
	}
	return appId
}
