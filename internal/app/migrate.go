package app

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// RunMigrations выполняет миграции базы данных
func RunMigrations(db *sqlx.DB) error {
	// Получаем *sql.DB из *sqlx.DB
	sqlDB := db.DB

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create postgres driver: %w", err)
	}

	// Путь к миграциям
	migrationsPath := "file://internal/migrations"

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	// Применяем миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("📋 No migrations to apply")
	} else {
		log.Println("✅ Migrations applied successfully")
	}

	// Получаем текущую версию
	version, dirty, err := m.Version()
	if err != nil {
		log.Printf("⚠️  Warning: Could not get migration version: %v", err)
	} else {
		log.Printf("📋 Current migration version: %d (dirty: %t)", version, dirty)
	}

	return nil
}
