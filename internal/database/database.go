package database

import (
    "database/sql"
    "fmt"
    "vend_erp/internal/config"

    _ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }

    // Настройка пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * 60) // 5 minutes

    return db, nil
}

func HealthCheck(db *sql.DB) error {
    return db.Ping()
}

func GetDBStats(db *sql.DB) sql.DBStats {
    return db.Stats()
}