package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"

	"order-service/internal/interfaces"
)

type Consumer struct {
	reader  *kafka.Reader
	service interfaces.OrderService
}

func NewConsumer(brokers []string, topic, groupID string, service interfaces.OrderService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafka.LastOffset,
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Println("Starting Kafka consumer...")

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer...")
			return c.reader.Close()
		default:
			// Чтение сообщения с таймаутом
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			// обработка сообщения
			if err := c.processMessage(msg); err != nil {
				log.Printf("Error processing message: %v", err)
				continue
			}
		}
	}
}

// processMessage обрабатывает полученное сообщение из Kafka.
func (c *Consumer) processMessage(msg kafka.Message) error {
	log.Printf("Received message: offset=%d, key=%s", msg.Offset, string(msg.Key))

	return c.service.ProcessOrder(msg.Value)
}
