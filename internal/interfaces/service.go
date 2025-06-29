package interfaces

import "order-service/internal/models"

type OrderService interface {
	ProcessOrder(data []byte) error
	GetOrder(orderUID string) (*models.Order, error)
	LoadCacheFromDB() error
	GetCacheMetrics() CacheMetrics
	GetCacheSize() int
}
