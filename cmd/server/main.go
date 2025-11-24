package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "vend_erp/internal/config"
    "vend_erp/internal/handlers"

    _ "github.com/lib/pq"
)

func main() {
    // Load environment variables from .env file
    if err := config.LoadEnv(); err != nil {
        log.Printf("Warning: Could not load .env file: %v", err)
    }
    
    cfg := config.LoadConfig()
    
    // Подключение к БД
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer db.Close()
    
    // Проверка подключения
    err = db.Ping()
    if err != nil {
        log.Fatal("Database ping failed:", err)
    }
    
    fmt.Println("✅ Successfully connected to database")
    
    mux := http.NewServeMux()
    handlers.SetupRoutes(mux, db)
    
    fmt.Printf("🚀 Server starting on %s\n", cfg.ServerAddress)
    fmt.Println("📊 Available routes:")
    fmt.Println("  GET  /dashboard")
    fmt.Println("  GET  /users")
    fmt.Println("  GET  /machines") 
    fmt.Println("  GET  /locations")
    fmt.Println("  GET  /api/stats")
    
    log.Fatal(http.ListenAndServe(cfg.ServerAddress, mux))
}