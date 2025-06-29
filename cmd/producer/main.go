package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Order struct {
	OrderUID        string    `json:"order_uid"`
	TrackNumber     string    `json:"track_number"`
	Entry           string    `json:"entry"`
	Delivery        Delivery  `json:"delivery"`
	Payment         Payment   `json:"payment"`
	Items           []Item    `json:"items"`
	Locale          string    `json:"locale"`
	InternalSig     string    `json:"internal_signature"`
	CustomerID      string    `json:"customer_id"`
	DeliveryService string    `json:"delivery_service"`
	ShardKey        string    `json:"shardkey"`
	SMID            int       `json:"sm_id"`
	DateCreated     time.Time `json:"date_created"`
	OofShard        string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	// Конфигурация Kafka
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	// Создание продюсера
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Тестовый заказ
	order := Order{
		OrderUID:    fmt.Sprintf("test_order_%d", time.Now().Unix()),
		TrackNumber: "TEST_TRACK_001",
		Entry:       "TEST",
		Delivery: Delivery{
			Name:    "Test User",
			Phone:   "+1234567890",
			Zip:     "12345",
			City:    "Test City",
			Address: "Test Address 123",
			Region:  "Test Region",
			Email:   "test@example.com",
		},
		Payment: Payment{
			Transaction:  fmt.Sprintf("test_tx_%d", time.Now().Unix()),
			Currency:     "USD",
			Provider:     "test_provider",
			Amount:       1000,
			PaymentDT:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []Item{
			{
				ChrtID:      12345,
				TrackNumber: "TEST_TRACK_001",
				Price:       500,
				RID:         "test_rid_001",
				Name:        "Test Product",
				Sale:        10,
				Size:        "M",
				TotalPrice:  450,
				NMID:        67890,
				Brand:       "Test Brand",
				Status:      200,
			},
		},
		Locale:          "en",
		InternalSig:     "test_sig",
		CustomerID:      "test_customer",
		DeliveryService: "test_delivery",
		ShardKey:        "1",
		SMID:            1,
		DateCreated:     time.Now(),
		OofShard:        "1",
	}

	// Конвертация в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %v", err)
	}

	// Отправка сообщения
	message := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.StringEncoder(orderJSON),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent successfully!\n")
	fmt.Printf("Order ID: %s\n", order.OrderUID)
	fmt.Printf("Partition: %d, Offset: %d\n", partition, offset)
	fmt.Printf("Test with: curl http://localhost:8081/order/%s\n", order.OrderUID)
}
