package database

import (
    "database/sql"
    "fmt"
    "vend_erp/config" // Changed from internal/config to config
    _ "github.com/jackc/pgx/v4/stdlib"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

    // Use "pgx" as driver name
    db, err := sql.Open("pgx", connStr)
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