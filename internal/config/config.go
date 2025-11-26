package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
)

// Config содержит конфигурацию приложения
type Config struct {
    DB *sql.DB
}

// NewConfig создает новую конфигурацию
func NewConfig() *Config {
    return &Config{}
}

// InitDB инициализирует подключение к базе данных
func (c *Config) InitDB() error {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "host=localhost port=5432 user=postgres password=postgres dbname=vend_erp sslmode=disable"
    }

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    if err := db.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    c.DB = db
    log.Println("✅ Database connection established")
    return nil
}

// GetDB возвращает подключение к базе данных
func (c *Config) GetDB() *sql.DB {
    return c.DB
}

// Close закрывает подключения
func (c *Config) Close() {
    if c.DB != nil {
        c.DB.Close()
    }
}
