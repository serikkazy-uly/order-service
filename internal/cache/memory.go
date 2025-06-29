package cache

import (
	"order-service/internal/interfaces"
	"order-service/internal/models"
	"sync"
)

type MemoryCache struct {
	mu      sync.RWMutex
	orders  map[string]*models.Order
	metrics interfaces.CacheMetrics
}

// Проверка соответствия интерфейсу
var _ interfaces.Cache = (*MemoryCache)(nil)

func NewMemoryCache() interfaces.Cache {
	return &MemoryCache{
		orders: make(map[string]*models.Order),
	}
}

func (c *MemoryCache) Set(orderUID string, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	orderCopy := *order
	c.orders[orderUID] = &orderCopy
}

func (c *MemoryCache) Get(orderUID string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	order, exists := c.orders[orderUID]
	if exists {
		c.metrics.Hits++
		orderCopy := *order
		return &orderCopy, true
	}

	c.metrics.Misses++

	return nil, false
}

func (c *MemoryCache) LoadFromDB(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		orderCopy := order
		c.orders[order.OrderUID] = &orderCopy
	}
}

func (c *MemoryCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.orders)
}

func (c *MemoryCache) GetMetrics() interfaces.CacheMetrics {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics
}
