package config

import "os"

type Config struct {
	Database DatabaseConfig
	Kafka    KafkaConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type ServerConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "database"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "orders"),
			GroupID: getEnv("KAFKA_GROUP_ID", "order-service"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8081"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
