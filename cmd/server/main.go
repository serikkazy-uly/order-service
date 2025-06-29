package main

import (
	"log"

	"order-service/internal/app"

	"github.com/joho/godotenv"
)

func main() {

	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf(" Warning: .env file not found: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}
	// Создаем приложение
	application := app.New()

	// Инициализируем
	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Запускаем
	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run application: %v", err)
	}

	// Корректно завершаем
	if err := application.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown application: %v", err)
	}
}
