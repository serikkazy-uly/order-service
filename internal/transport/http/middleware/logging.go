package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware - это middleware для логирования HTTP-запросов.
// Логирует метод запроса, URL и время выполнения запроса
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызываем следующий handler
		next.ServeHTTP(w, r)

		// Логируем после выполнения
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}
