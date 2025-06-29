package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/interfaces"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/internal/transport/http"
	"order-service/internal/transport/http/handlers"
	"order-service/internal/transport/kafka"
)

// App представляет основное приложение
type App struct {
	config        *config.Config
	db            *sqlx.DB
	cache         interfaces.Cache
	service       interfaces.OrderService
	httpServer    *http.Server
	kafkaConsumer *kafka.Consumer
}

// New создает новый экземпляр приложения
func New() *App {
	return &App{}
}

// Initialize инициализирует все компоненты приложения
func (a *App) Initialize() error {
	log.Println("Initializing Order Service...")

	// 1. Загружаем конфигурацию
	a.config = config.Load()
	log.Printf("Config loaded: DB=%s:%s, Kafka=%v, Port=%s",
		a.config.Database.Host, a.config.Database.Port,
		a.config.Kafka.Brokers, a.config.Server.Port)

	// 2. Подключаемся к базе данных
	db, err := ConnectDatabase(a.config.Database)
	if err != nil {
		return err
	}
	a.db = db
	log.Println("Database connected successfully")

	// 3. Запускаем миграции
	log.Println("Running database migrations...")
	if err := RunMigrations(a.db); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// 4. Создаем слои приложения
	a.cache = cache.NewMemoryCache()
	repo := repository.NewOrderRepository(a.db)
	a.service = service.NewOrderService(repo, a.cache)

	// 5. Загружаем кеш из БД
	if err := a.loadCache(); err != nil {
		log.Printf("Warning: Failed to load cache: %v", err)
	}

	// 6. Инициализируем HTTP сервер
	a.initHTTPServer()

	// 7. Инициализируем Kafka consumer
	a.initKafkaConsumer()

	log.Println("Application initialized successfully")
	return nil
}

// Run запускает приложение
func (a *App) Run() error {
	log.Println("Starting Order Service...")

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka consumer
	go func() {
		log.Println("Starting Kafka consumer...")
		if err := a.kafkaConsumer.Start(ctx); err != nil {
			log.Printf("Kafka consumer error: %v", err)
		}
	}()

	// Запускаем HTTP сервер
	go func() {
		log.Printf("HTTP server starting on port %s", a.config.Server.Port)
		log.Printf("Web interface: http://localhost:%s", a.config.Server.Port)
		log.Printf("API: http://localhost:%s/health", a.config.Server.Port)
		log.Printf("Order API: http://localhost:%s/order/{order_uid}", a.config.Server.Port)

		if err := a.httpServer.Start(); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Ожидаем сигнал для завершения
	return a.waitForShutdown(ctx, cancel)
}

// Shutdown корректно завершает работу приложения
func (a *App) Shutdown() error {
	log.Println("Shutting down application...")

	if a.db != nil {
		a.db.Close()
		log.Println("Database connection closed")
	}

	log.Println("Application shutdown completed")
	return nil
}

// loadCache загружает кеш из базы данных
func (a *App) loadCache() error {
	if err := a.service.LoadCacheFromDB(); err != nil {
		return err
	}

	size := a.service.GetCacheSize()
	log.Printf("Loaded %d orders into cache", size)
	return nil
}

// initHTTPServer инициализирует HTTP сервер
func (a *App) initHTTPServer() {
	orderHandler := handlers.NewOrderHandler(a.service)
	a.httpServer = http.NewServer(a.config.Server.Port, orderHandler)
}

// initKafkaConsumer инициализирует Kafka consumer
func (a *App) initKafkaConsumer() {
	a.kafkaConsumer = kafka.NewConsumer(
		a.config.Kafka.Brokers,
		a.config.Kafka.Topic,
		a.config.Kafka.GroupID,
		a.service,
	)
}

// waitForShutdown ожидает сигнал для завершения работы
func (a *App) waitForShutdown(_ context.Context, cancel context.CancelFunc) error {
	// Канал для получения сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал
	<-quit
	log.Println("Received shutdown signal")

	// Отменяем контекст для остановки Kafka consumer
	cancel()

	// Останавливаем HTTP сервер
	if a.httpServer != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Println("HTTP server stopped")
	}

	return nil
}
