package database

import (
	"aurma_product/internal/config"
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	maxAttempts       = 10
	connectionTimeout = 5 * time.Second
)

// Connect устанавливает соединение с MySQL базой данных.
func Connect(cfg *config.Config) (*sqlx.DB, error) {
	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("Attempting to connect to MySQL (attempt %d/%d)", attempt, maxAttempts)

		db, err := sqlx.Open("mysql", databaseURL)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			time.Sleep(connectionTimeout)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
		err = db.PingContext(ctx)
		cancel()

		if err == nil {
			return db, nil
		}
		log.Printf("Failed to ping database: %v", err)
		db.Close()
		time.Sleep(connectionTimeout)
	}

	return nil, fmt.Errorf("failed to connect to MySQL database after %d attempts", maxAttempts)
}
