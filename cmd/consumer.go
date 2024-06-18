package cmd

import (
	"backend-exercise/utils"
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var kafkaCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume messages from Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		consumeKafkaMessages()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
}

func consumeKafkaMessages() {
	r := utils.InitKafkaReader()
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message from Kafka: %v", err)
			continue
		}
		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
