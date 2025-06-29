# Загружаем переменные из .env файла если он существует
# ifneq (,$(wildcard ./.env))
#     include .env
#     export
# endif

# Параметры по умолчанию (только для разработки)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= default_user
DB_NAME ?= default_name
DB_PASSWORD ?= default_password

DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATIONS_PATH = internal/migrations
DOCKER_COMPOSE_FILE = docker-compose.yml

.PHONY: help docker-up docker-down docker-status run build deps test clean migrate-up migrate-down migrate-version setup

# Показать справку
help:
	@echo "Order Service - Команды Make"
	@echo ""
	@echo " Docker команды:"
	@echo "  make docker-up       # Запустить PostgreSQL + Kafka"
	@echo "  make docker-down     # Остановить контейнеры"
	@echo "  make docker-status   # Статус контейнеров"
	@echo ""
	@echo " Приложение:"
	@echo "  make run            # Запустить сервис (автомиграции включены)"
	@echo "  make build          # Собрать приложение"
	@echo "  make test           # Запустить тесты"
	@echo ""
	@echo "  Миграции (опционально - встроены в приложение):"
	@echo "  make migrate-up     # Применить миграции вручную"
	@echo "  make migrate-down   # Откатить миграции"
	@echo "  make migrate-version # Показать версию миграций"
	@echo ""
	@echo "  Утилиты:"
	@echo "  make deps           # Установить зависимости"
	@echo "  make setup          # Полная настройка проекта"
	@echo "  make clean          # Очистить данные"

# Запуск docker-compose с автозагрузкой .env
docker-up:
	@echo " Запуск Docker контейнеров..."
	@if [ -f .env ]; then \
		docker compose --env-file .env -f $(DOCKER_COMPOSE_FILE) up -d; \
	else \
		echo " .env файл не найден, используем значения по умолчанию"; \
		docker compose -f $(DOCKER_COMPOSE_FILE) up -d; \
	fi
	@echo "Контейнеры запущены"

# Остановить docker-compose
docker-down:
	@echo "Остановка Docker контейнеров..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "Контейнеры остановлены"

# Статус контейнеров
docker-status:
	@echo " Статус Docker контейнеров:"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Запустить приложение (миграции применяются автоматически)
run:
	@echo " Запуск Order Service..."
	@echo " Миграции БД применяются автоматически при старте"
	go run cmd/server/main.go

# Собрать приложение
build:
	@echo " Сборка приложения..."
	@mkdir -p bin
	go build -o bin/server cmd/server/main.go
	@echo "Приложение собрано: bin/server"


# Установить зависимости
deps:
	@echo "Установка зависимостей..."
	go mod download
	go mod tidy
	@echo " Установка дополнительных пакетов для миграций..."
	go get github.com/golang-migrate/migrate/v4
	go get github.com/golang-migrate/migrate/v4/database/postgres
	go get github.com/golang-migrate/migrate/v4/source/file
	go get github.com/joho/godotenv
	@echo " Зависимости установлены"


# Показать версию миграций
migrate-version:
	@echo " Версия миграций:"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version; \
	else \
		echo " golang-migrate не установлен"; \
		echo " Миграции встроены в приложение"; \
	fi

# Полная настройка проекта
setup: deps docker-up
	@echo " Ожидание готовности сервисов..."
	@sleep 15
	@echo " Проверка подключения к базе данных..."
	@docker exec -it orders_postgres psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT current_database();" || true
	@echo ""
	@echo " Настройка завершена!"
	@echo " Теперь запустите: make run"

# Очистка данных
clean:
	@echo " Очистка данных..."
	make docker-down
	@echo "  Удаление данных PostgreSQL..."
	docker volume rm order-service_postgres_data 2>/dev/null || true
	@echo " Очистка бинарных файлов..."
	rm -rf bin/
	@echo " Очистка завершена"


# Логи контейнеров
logs:
	@echo " Логи Docker контейнеров:"
	@echo "--- PostgreSQL ---"
	@docker logs orders_postgres --tail 10 2>/dev/null || echo "PostgreSQL не запущен"
	@echo "--- Kafka ---"
	@docker logs orders_kafka --tail 10 2>/dev/null || echo "Kafka не запущен"

# Проверка здоровья системы
health:
	@echo " Проверка здоровья системы:"
	@echo " Docker контейнеры:"
	@docker ps --format "table {{.Names}}\t{{.Status}}" | grep orders || echo "Контейнеры не запущены"
	@echo ""
	@echo "  База данных:"
	@docker exec orders_postgres pg_isready -U $(DB_USER) -d $(DB_NAME) 2>/dev/null && echo " PostgreSQL готов" || echo " PostgreSQL недоступен"
	@echo ""
	@echo " Приложение:"
	@curl -s http://localhost:8081/health >/dev/null 2>&1 && echo " Order Service работает" || echo " Order Service недоступен"

