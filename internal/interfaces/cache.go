package interfaces

import "order-service/internal/models"

type Cache interface {
	Set(orderUID string, order *models.Order)
	Get(orderUID string) (*models.Order, bool)
	LoadFromDB(orders []models.Order)
	Size() int
	GetMetrics() CacheMetrics
}

type CacheMetrics struct {
	Hits   int64
	Misses int64
}
