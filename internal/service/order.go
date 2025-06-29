package service

import (
	"encoding/json"
	"fmt"
	"log"

	"order-service/internal/interfaces"

	"order-service/internal/models"
)

type orderService struct {
	repo  interfaces.OrderRepository
	cache interfaces.Cache
}

// Проверка соответствия интерфейсу
var _ interfaces.OrderService = (*orderService)(nil)

func NewOrderService(r interfaces.OrderRepository, c interfaces.Cache) interfaces.OrderService {
	return &orderService{
		repo:  r,
		cache: c,
	}
}

// обработка заказа из Kafka
func (s *orderService) ProcessOrder(data []byte) error {
	// Парсинг JSON
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	// Валидация данных
	if err := s.validateOrder(&order); err != nil {
		log.Printf("Invalid order data: %v", err)
		return fmt.Errorf("invalid order data: %w", err)
	}

	// Сохранение в БД
	if err := s.repo.CreateOrder(&order); err != nil {
		return fmt.Errorf("failed to save order to database: %w", err)
	}

	//  Обновление кеша
	s.cache.Set(order.OrderUID, &order)

	log.Printf("Order %s processed successfully", order.OrderUID)
	return nil
}

// получение заказа (кеш + БД)
func (s *orderService) GetOrder(orderUID string) (*models.Order, error) {
	// Проверяем кеш
	if order, exists := s.cache.Get(orderUID); exists {
		log.Printf("Cache hit for order: %s", orderUID)
		return order, nil
	}

	// Если в кеше нет, обращаемся к БД
	log.Printf("Cache miss for order: %s, fetching from database", orderUID)

	order, err := s.repo.GetOrder(orderUID)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	s.cache.Set(orderUID, order)

	return order, nil
}

// восстановление кеша при старте
func (s *orderService) LoadCacheFromDB() error {
	log.Println("Loading cache from database...")

	orders, err := s.repo.GetAllOrders()
	if err != nil {
		return fmt.Errorf("failed to load orders from database: %w", err)
	}

	s.cache.LoadFromDB(orders)
	log.Printf("Loaded %d orders into cache", len(orders))

	return nil
}

// Валидация заказа
func (s *orderService) validateOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid is required")
	}

	if order.TrackNumber == "" {
		return fmt.Errorf("track_number is required")
	}

	if len(order.Items) == 0 {
		return fmt.Errorf("order must contain at least one item")
	}

	return nil
}

func (s *orderService) GetCacheMetrics() interfaces.CacheMetrics {
	return s.cache.GetMetrics()
}

func (s *orderService) GetCacheSize() int {
	return s.cache.Size()
}
