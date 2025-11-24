package main

import (
    "fmt"
    "log"
    "net/http"
    
    // Import your local packages
    "vend_erp/config"
    "vend_erp/internal/handlers"
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

    // Setup routes using handlers
    mux := http.NewServeMux()
    
    // Setup your existing routes
    handlers.SetupRoutes(mux, db)

    // Add a simple health check route
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        // Check database connection
        err := db.Ping()
        if err != nil {
            http.Error(w, "Database connection failed", http.StatusInternalServerError)
            return
        }
        
        // Check if migrations table exists (indicates successful migrations)
        var migrationCount int
        err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&migrationCount)
        if err != nil {
            http.Error(w, "Migrations not properly applied", http.StatusInternalServerError)
            return
        }
        
        fmt.Fprintf(w, `{
            "status": "healthy",
            "database": "connected",
            "migrations_applied": %d,
            "database_name": "%s"
        }`, migrationCount, cfg.DBName)
    })

    // Start server
    port := ":8080"
    log.Printf("🚀 Vend ERP Server starting on http://localhost%s", port)
    log.Printf("📊 Database: %s@%s:%d/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
    log.Printf("🗃️  Migrations applied from: %s", migrationsPath)
    log.Printf("🔧 Health check: http://localhost%s/health", port)
    
    if err := http.ListenAndServe(port, mux); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}