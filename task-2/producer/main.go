package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// User is a simple record example
type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
	FavoriteColor  string `json:"favorite_color"`
}

func main() {
	bootstrapServers := "localhost:9093" // Replace with your Kafka broker address
	topic := "topic-2"                   // Replace with your Kafka topic

	// SSL and SASL configuration
	sslConfig := &kafka.ConfigMap{
		"bootstrap.servers":                   bootstrapServers, // Replace with your Kafka broker address
		"security.protocol":                   "SASL_SSL",
		"ssl.ca.location":                     "../ca.crt", // Path to your CA certificate
		"ssl.keystore.location":               "../kafka-1-creds/kafka.kafka-1.keystore.pkcs12",
		"ssl.keystore.password":               "your-password",
		"ssl.key.password":                    "your-password", // Password for the client key
		"sasl.mechanism":                      "PLAIN",
		"sasl.username":                       "producer",
		"sasl.password":                       "pass",
		"enable.ssl.certificate.verification": false,
	}

	// Create a new producer
	p, err := kafka.NewProducer(sslConfig)
	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	defer p.Close()

	fmt.Printf("Created Producer %v\n", p)

	// Prepare the message payload
	value := User{
		Name:           "First user",
		FavoriteNumber: 42,
		FavoriteColor:  "blue",
	}
	payload, err := json.Marshal(value)
	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	// Delivery channel for message delivery reports
	deliveryChan := make(chan kafka.Event)

	// Produce the message
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, deliveryChan)
	if err != nil {
		fmt.Printf("Produce failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
}
