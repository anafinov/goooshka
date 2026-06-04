package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"student-cert-service/internal/domain"
	"student-cert-service/internal/repository"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is required")
	}

	time.Sleep(2 * time.Second)

	repo, err := repository.NewPostgresRepository(dbURL)
	if err != nil {
		log.Fatalf("failed to init repository: %v", err)
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"cert_requests",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var req domain.Request
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			log.Printf("Received request ID %d. Processing...", req.ID)

			time.Sleep(5 * time.Second)

			if err := repo.UpdateRequestStatus(req.ID, "ready"); err != nil {
				log.Printf("Error updating request status: %v", err)
			} else {
				log.Printf("Request ID %d is ready", req.ID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
