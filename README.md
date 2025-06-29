# Order Service (Сервис заказов) - Level 0

## Техническое задание

Демонстрационный сервис с простейшим интерфейсом для отображения данных о заказах с использованием Kafka, PostgreSQL и кеширования.

---

## Цель проекта

Создание микросервиса на Go, который:
- Получает данные заказов из очереди сообщений (Kafka)
- Сохраняет их в базу данных (PostgreSQL) 
- Кеширует в памяти для быстрого доступа
- Предоставляет HTTP API и веб-интерфейс для просмотра данных

---

## Технический стек: 
`Go 1.21+, PostgreSQL 15+, Kafka, In-Memory (Go map), Docker, Docker Compose, Gorilla Mux, HTML/CSS/JS` 


## Функциональные требования

### 1. База данных

- **Развернуть PostgreSQL локально**
  - Создать БД `orders_db`
  - Настроить пользователя `orders_user` с правами доступа
  - Спроектировать схему для хранения данных заказов

### 2. Обработка данных

- **Kafka Consumer**
  - Подключиться к брокеру сообщений
  - Подписаться на топик `orders`
  - Обрабатывать сообщения в реальном времени

- **Валидация данных**
  - Проверять корректность структуры сообщений
  - Логировать некорректные данные
  - Игнорировать невалидные сообщения

### 3. Хранение данных

- **PostgreSQL интеграция**
  - Парсить JSON сообщения
  - Сохранять данные в реляционные таблицы
  - Использовать транзакции для целостности данных

- **Кеширование**
  - Хранить данные заказов в памяти (Go map)
  - Обновлять кеш при получении новых данных
  - Восстанавливать кеш из БД при перезапуске сервиса

### 4. HTTP API

- **REST endpoints**
  ```
  GET /order/{order_uid} - получение данных заказа
  GET /health           - проверка состояния сервиса
  ```

- **Веб-интерфейс**
  - HTML страница с формой ввода order_uid
  - JavaScript для взаимодействия с API
  - Отображение данных в удобном формате

---

## Модель данных

### Структура заказа (JSON)
```json
{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK", 
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}
```

### Компоненты системы

| Компонент | Ответственность |
|:----------|:----------------|
| **Kafka Consumer** | Получение и обработка сообщений из Kafka |
| **Service Layer** | Бизнес-логика валидации и обработки заказов |
| **Repository Layer** | Взаимодействие с PostgreSQL |
| **Cache Layer** | Управление кешем в памяти |
| **HTTP Handlers** | Обработка HTTP запросов |
| **Web Interface** | Пользовательский интерфейс |

---

## Сценарии тестирования

### 1. Базовый сценарий [Данный проект запускается ч/з Docker и .env (см ниже)]
```bash
# 1. Запустить сервис
make run

# 2. Отправить сообщение в Kafka
go run cmd/producer/main.go

# 3. Проверить API
curl http://localhost:8081/order/b563feb7b2b84b6test

# 4. Открыть веб-интерфейс и ввести order_uid
open http://localhost:8081
```

### 2. Тест производительности кеша
```bash
# 1. Отправить заказ
go run cmd/producer/main.go

# 2. Первый запрос (из БД)
time curl http://localhost:8081/order/b563feb7b2b84b6test

# 3. Второй запрос (из кеша) - должен быть быстрее
time curl http://localhost:8081/order/b563feb7b2b84b6test
```

### 3. Тест восстановления кеша
```bash
# 1. Отправить заказ и убедиться, что он в кеше
# 2. Перезапустить сервис
# 3. Проверить, что данные доступны без задержки
```

### `GET /order/{order_uid}`
**Описание:** Получение информации о заказе по идентификатору.

**Параметры:**
- `order_uid` *(path, string, required)* - Уникальный идентификатор заказа

**Пример запроса:**
```bash
curl -X GET http://localhost:8081/order/b563feb7b2b84b6test
```

**Успешный ответ (200 OK):**
```json
{
  "order_uid": "b563feb7b2b84b6test",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL",
  "delivery": {
    "name": "Test Testov",
    "phone": "+9720000000",
    "zip": "2639809",
    "city": "Kiryat Mozkin",
    "address": "Ploshad Mira 15",
    "region": "Kraiot",
    "email": "test@gmail.com"
  },
  "payment": {
    "transaction": "b563feb7b2b84b6test",
    "request_id": "",
    "currency": "USD",
    "provider": "wbpay",
    "amount": 1817,
    "payment_dt": 1637907727,
    "bank": "alpha",
    "delivery_cost": 1500,
    "goods_total": 317,
    "custom_fee": 0
  },
  "items": [
    {
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }
  ],
  "locale": "en",
  "internal_signature": "",
  "customer_id": "test",
  "delivery_service": "meest",
  "shardkey": "9",
  "sm_id": 99,
  "date_created": "2021-11-26T06:22:19Z",
  "oof_shard": "1"
}
```

**Ошибка - заказ не найден (404 NOT FOUND):**
```json
{
  "error": "Order not found",
  "order_uid": "invalid_uid"
}
```

**Ошибка сервера (500 INTERNAL SERVER ERROR):**
```json
{
  "error": "Internal server error",
  "message": "Database connection failed"
}
```

### `GET /health`
**Описание:** Проверка состояния сервиса и его компонентов.

**Пример запроса:**
```bash
curl -X GET http://localhost:8081/health
```

**Успешный ответ (200 OK):**
```json
{
  "status": "ok",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "database": {
      "status": "connected",
      "response_time_ms": 5
    },
    "kafka": {
      "status": "connected",
      "topic": "orders"
    },
    "cache": {
      "status": "active",
      "size": 42,
      "hit_rate": 0.85
    }
  }
}
```

**Ошибка - проблемы с сервисами (503 SERVICE UNAVAILABLE):**
```json
{
  "status": "degraded",
  "timestamp": "2024-01-15T10:30:00Z",
  "services": {
    "database": {
      "status": "error",
      "error": "Connection timeout"
    },
    "kafka": {
      "status": "connected",
      "topic": "orders"
    },
    "cache": {
      "status": "active",
      "size": 42,
      "hit_rate": 0.85
    }
  }
}
```
---

## Мониторинг и логирование

### Логи
Сервис логирует следующие события:
- Получение сообщений из Kafka
- Сохранение данных в БД
- Обращения к кешу
- HTTP запросы
- Ошибки обработки
---

## Известные ограничения

1. **Кеш в памяти** - данные теряются при перезапуске (восстанавливаются из БД)
2. **Односерверная архитектура** - нет горизонтального масштабирования
3. **Простая аутентификация** - нет защиты API
4. **Базовая обработка ошибок** - minimal retry logic
---

##  Демонстрация работы

Для сдачи проекта необходимо записать короткое видео (3-5 минут), демонстрирующее:

1. **Запуск инфраструктуры** - `make docker-up`
2. **Запуск сервиса** - `make run`
3. **Отправка сообщения** - `go run cmd/producer/main.go`
4. **Веб-интерфейс** - ввод order_uid и получение данных
5. **API тестирование** - curl запрос к `/order/{id}`
6. **Перезапуск сервиса** - демонстрация восстановления кеша

### Что показать в видео:
- Рабочий веб-интерфейс
- Корректное отображение JSON данных  
- Скорость работы кеша (повторные запросы)
- Обработка несуществующих заказов
- Логи в консоли сервиса
---


## Поддержка
При возникновении вопросов и т.п.:
- Создай Issue

## Для ревьювера!

 Почему sqlx?
 - Автоматический маппинг по db тегам
 - Именованные параметры - понятно и безопасно
_, err = db.NamedExec (`Автоматически берет поля из структуры!`)
- Минимальный API поверх стандартной библиотеки
и не нужно изучать новую парадигму (как в GORM)
-  видишь SQL запросы
---

# Быстрый старт проекта

## 1. Подготовка окружения

**Предварительные требования:** смотри выше [Технический стек]

### Клонировать репозиторий
```bash
git clone https://github.com/serikkazy-uly/order-service
cd order-service
```

### Установить зависимости
```bash
go mod download
go mod tidy
```

## 2. Настройка переменных окружения

### Создать .env файл
```bash
# Создать .env файл
cat > .env << EOF
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=dk_orders_user
DB_PASSWORD=dk_orders_pswd
DB_NAME=orders_postgres

# Kafka configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=orders
KAFKA_GROUP_ID=order-service-group

# Server configuration
SERVER_PORT=8081
EOF
```

## 3. Запуск инфраструктуры

### Запуск PostgreSQL и Kafka через Docker
```bash
make docker-up
```

### Проверка что контейнеры запущены
```bash
docker ps
# Должны быть: orders_postgres, orders_kafka, orders_zookeeper (все в статусе Up)
```

### Ожидание готовности сервисов
```bash
# Подождать 30-60 секунд пока PostgreSQL и Kafka полностью запустятся
sleep 30
```

## 4. Установка golang-migrate (если еще не установлен)

### macOS:
```bash
brew install golang-migrate
```

### Linux:
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### Через Go:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Установка зависимостей для миграций (если не установлены)
```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
go get github.com/joho/godotenv  # Для загрузки .env файла
```

## 5. Применение миграций базы данных автоматически

### Проверка миграций (опционально)
```bash
# Проверка версии миграций
make migrate-version

# Проверка таблиц в БД
docker exec -it orders_postgres psql -U dk_orders_user -d dk_orders_db -c "\dt"
```

## 6. Сборка сервиса (опционально - для production)

```bash
# Сборка бинарного файла (опционально)
make build
```

**Примечание:** Этот шаг **опционален** для разработки, так как `make run` компилирует и запускает сразу.

## 7. Запуск приложения

### Старт сервера (включает автоматические миграции)
```bash
make run
```

**Ожидаемый вывод:**
```
Initializing Order Service...
Config loaded: DB=localhost:5432, Kafka=[localhost:9092], Port=8081
Database connected successfully
Running database migrations...
Migrations applied successfully
Loaded 1 orders into cache
Application initialized successfully
Starting Order Service...
Starting Kafka consumer...
HTTP server starting on port 8081
Web interface: http://localhost:8081
API: http://localhost:8081/health
```

## 8. Проверка работы системы

### Health check
```bash
curl http://localhost:8081/health
```

### Получить заказ по ID
```bash
curl http://localhost:8081/order/b563feb7b2b84b6test
```

### Получить заказ с метриками времени
```bash
curl -w "\nTime: %{time_total}s\n" http://localhost:8081/order/b563feb7b2b84b6test
```

### Веб-интерфейс
```bash
# Открыть в браузере:
open http://localhost:8081/
# или просто перейти на http://localhost:8081/
```

**Действия в браузере:**
1. Введите Order ID: `b563feb7b2b84b6test`
2. Нажмите "Search"
3. Проверьте отображение JSON данных

## 9. Тестирование Kafka Producer (опционально)

### Создать тестовый продюсер (если не существует)
```bash
# Создать файл cmd/producer/main.go с кодом продюсера
```

### Запустить продюсер
```bash
go run cmd/producer/main.go
```

## `!!!` Критические моменты последовательности:

1. **Docker ОБЯЗАТЕЛЬНО до миграций** - БД должна быть запущена
2. **Ожидание готовности сервисов** - PostgreSQL нужно время на инициализацию
3. **Миграции до запуска APP** - таблицы должны существовать
4. **Сборка опциональна** - `make run` компилирует автоматически
5. **Проверки после запуска APP** - сервис должен полностью стартовать

## 🔧 Решение проблем по шагам:

### Если шаг 3 не работает:
```bash
# Проверить Docker
docker --version
docker compose --version

# Пересоздать контейнеры
make docker-down
make docker-up
```

### Если шаг 5 не работает:
```bash
# Проверить что migrate установлен
migrate -version

# Проверить подключение к БД
docker exec -it orders_postgres psql -U dk_orders_user -d dk_orders_db -c "SELECT 1;"
```

### Если шаг 7 не работает:
```bash
# Проверить .env файл
cat .env

# Проверить что .env загружается в docker-compose через заданные переменные
docker compose --env-file .env -f docker-compose.yml config\n