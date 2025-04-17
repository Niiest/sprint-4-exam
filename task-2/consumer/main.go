package main

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	bootstrapServers := "localhost:9093" // Replace with your Kafka broker address
	topic := "topic-1"                   // Replace with your Kafka topic
	groupID := "my-consumer-group"       // Consumer group ID

	// SSL configuration
	sslConfig := &kafka.ConfigMap{
		"bootstrap.servers":        bootstrapServers,
		"security.protocol":        "SASL_SSL",
		"ssl.ca.location":          "../ca.crt",                    // Path to your CA certificate
		"ssl.certificate.location": "../kafka-1-creds/kafka-1.crt", // Path to your client certificate
		"ssl.truststore.location":  "../kafka-1-creds/kafka.kafka-1.truststore.jks",
		"ssl.truststore.password":  "../kafka-1-creds/your-password",
		"ssl.key.location":         "../kafka-1-creds/kafka-1.key", // Path to your client key
		"ssl.key.password":         "yandex",                       // Password for the client key
		"group.id":                 groupID,                        // Consumer group ID
		"auto.offset.reset":        "earliest",                     // Start reading at the earliest message
		"sasl.mechanism":           "PLAIN",
		"sasl.username":            "consumer",
		"sasl.password":            "pass",
	}

	// Create a new consumer
	c, err := kafka.NewConsumer(sslConfig)
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()

	// Subscribe to the topic
	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		fmt.Printf("Failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Consumer created and subscribed to topic %s\n", topic)

	// Poll for messages
	for {
		msg, err := c.ReadMessage(-1) // -1 means wait indefinitely for a message
		if err != nil {
			// Handle errors
			if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsFatal() {
				fmt.Printf("Fatal error: %v\n", kafkaErr)
				break
			}
			// Ignore non-fatal errors
			continue
		}

		// Process the message
		fmt.Printf("Received message: %s\n", string(msg.Value))
	}
}
