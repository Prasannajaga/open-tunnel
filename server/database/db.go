package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"opentunnel/server/config"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func TestConnection() {
	cfg := config.NewConfig()
	fmt.Println("cfg", cfg)
	if err := Connect(cfg); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer Close()
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
