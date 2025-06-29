package interfaces

import "order-service/internal/models"

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	GetOrder(orderUID string) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
}
