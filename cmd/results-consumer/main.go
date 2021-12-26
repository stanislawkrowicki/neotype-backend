package main

import (
	"log"
	"neotype-backend/pkg/rabbitmq"
	"neotype-backend/pkg/results"
)

const queueName = "results"

func main() {
	results.InitConsumer()
	rabbit := rabbitmq.New()
	err := rabbit.Connect(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to connect to rabbitmq: %s", err)
		return
	}

	messages, err := rabbit.Channel.Consume(
		rabbit.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register results consumer: %s", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for delivery := range messages {
			results.ConsumeResult(delivery.Body)
		}
	}()

	log.Printf(" [*] Listening for messages...")

	<-forever
}
