package main

import (
    "log"
    "net/http"
    
    "vend_erp/config"
    "vend_erp/migrations"

    _ "github.com/lib/pq"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Connect to database
    db, err := config.ConnectDB(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    log.Printf("Database connected successfully!")

    // Run migrations from migrations folder
    migrationsPath := "./migrations"
    if err := migrations.RunMigrations(db, migrationsPath); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // Setup routes using handlers package
    router := setupRoutes(db)

    // Start server
    port := ":8080"
    log.Printf("🚀 Vend ERP Server starting on http://localhost%s", port)
    log.Printf("📊 Database: %s@%s:%d/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
    log.Printf("🗃️  Migrations applied from: %s", migrationsPath)
    if err := http.ListenAndServe(port, router); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}