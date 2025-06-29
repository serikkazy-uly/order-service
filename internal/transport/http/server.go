package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"order-service/internal/transport/http/handlers"
	"order-service/internal/transport/http/middleware"
)

type Server struct {
	server *http.Server
}

func NewServer(port string, orderHandler *handlers.OrderHandler) *Server {
	r := mux.NewRouter()

	// Web pages
	r.HandleFunc("/health", orderHandler.Health).Methods("GET")
	r.HandleFunc("/order/{order_uid}", orderHandler.GetOrder).Methods("GET")

	// Главная страница
	r.HandleFunc("/", serveHome).Methods("GET")

	// Middleware
	r.Use(middleware.LoggingMiddleware) // Логируем запросы
	r.Use(middleware.CorsMiddleware)    // CORS policy (Разрешаем CORS JavaScript запрос)

	// Настройка CORS
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{server: srv}
}

// Start запускает HTTP сервер на заданном порту.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully останавливает сервер, позволяя завершить обработку текущих запросов.
// Это важно для предотвращения потери данных и обеспечения корректного завершения работы сервиса.
// Принимает контекст для возможности установки таймаута или отмены операции.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}
