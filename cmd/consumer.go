package cmd

import (
	"backend-exercise/utils"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	kafkaCmd = &cobra.Command{
		Use:   "consumer",
		Short: "Run the consumer to consume data",
		Run: func(cmd *cobra.Command, args []string) {
			consumeKafkaMessages()
		},
	}
	modelKeyPrefix string
)

func init() {
	kafkaCmd.Flags().StringVar(&modelKeyPrefix, "key-prefix", "", "Key prefix to filter messages")
	rootCmd.AddCommand(kafkaCmd)
}

func consumeKafkaMessages() {
	r := utils.InitKafkaReader(utils.GetKafkaTopic())
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message from Kafka: %v", err)
			continue
		}

		if modelKeyPrefix != "" && !hasKeyPrefix(m.Key, modelKeyPrefix) {
			continue
		}

		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}

func hasKeyPrefix(key []byte, prefix string) bool {
	return len(key) >= len(prefix) && string(key[:len(prefix)]) == prefix
}
