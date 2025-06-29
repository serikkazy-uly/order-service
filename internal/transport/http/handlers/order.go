package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"order-service/internal/interfaces"

	apperrors "order-service/internal/errors" // кастомнаые ошибки
)

type OrderHandler struct {
	service interfaces.OrderService
}

func NewOrderHandler(service interfaces.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// обработк GET /order/{order_uid}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметр из URL
	vars := mux.Vars(r)
	orderUID := vars["order_uid"]

	if orderUID == "" {
		h.writeError(w, "order_uid is required", http.StatusBadRequest)
		return
	}

	//  Получаем заказ
	order, err := h.service.GetOrder(orderUID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrOrderNotFound):
			h.writeError(w, "Order not found", http.StatusNotFound)
		case errors.Is(err, apperrors.ErrInvalidOrderUID):
			h.writeError(w, "Invalid order UID", http.StatusBadRequest)
		default:
			h.writeError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.writeJSON(w, order)
}

// Health check endpoint - Автоматическая проверка доступности :8081/health
func (h *OrderHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "order-service",
	}

	h.writeJSON(w, response)
}

// Служебные методы для HTTP ответов
func (h *OrderHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.writeError(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
