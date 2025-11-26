package main

import (
    "log"
    "net/http"

    "vend_erp/internal/database"
    "vend_erp/internal/handlers"
)

func main() {
    // Инициализация базы данных
    if err := database.Init(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer database.DB.Close()

    // Настройка маршрутов
    mux := http.NewServeMux()
    handlers.SetupRoutes(mux, database.DB)

    // Запуск сервера
    log.Println("🚀 Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
