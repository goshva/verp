#!/bin/bash

set -e

echo "🚀 Исправляем проблемы с импортами и структурой проекта..."

# Проверяем структуру проекта
echo "🔍 Проверяем структуру проекта..."
find . -name "*.go" -type f | head -10

# Проверяем наличие config пакета
echo "📁 Ищем config пакет..."
find . -name "config.go" -type f

# Создаем недостающий config пакет если его нет
if [ ! -f "internal/config/config.go" ]; then
    echo "📝 Создаем internal/config/config.go..."
    mkdir -p internal/config
    cat > internal/config/config.go << 'ENDOFFILE'
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
ENDOFFILE
    echo "✅ config.go создан"
fi

# Проверяем database.go
echo "🔧 Проверяем internal/database/database.go..."
if [ -f "internal/database/database.go" ]; then
    echo "📝 Обновляем импорты в database.go..."
    cat > internal/database/database.go << 'ENDOFFILE'
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

// DB глобальная переменная базы данных
var DB *sql.DB

// Init инициализирует подключение к базе данных
func Init() error {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        connStr = "host=localhost port=5432 user=postgres password=postgres dbname=vend_erp sslmode=disable"
    }

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open database: %v", err)
    }

    if err := DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %v", err)
    }

    log.Println("✅ Database connection established")
    return nil
}

// GetDB возвращает подключение к базе данных
func GetDB() *sql.DB {
    return DB
}
ENDOFFILE
    echo "✅ database.go обновлен"
fi

# Обновляем импорты в других файлах
echo "🔄 Обновляем импорты в handlers..."

# Создаем временный router.go с правильными импортами
if [ -f "internal/handlers/router.go" ]; then
    echo "📝 Обновляем router.go..."
    cat > internal/handlers/router.go << 'ENDOFFILE'
package handlers

import (
    "database/sql"
    "fmt"
    "html/template"
    "net/http"
    "path/filepath"
)

var tmpl *template.Template

func init() {
    tmpl = template.New("")
    
    // Загружаем шаблоны из всех поддиректорий
    patterns := []string{
        "internal/templates/*.html",
        "internal/templates/layouts/*.html",
        "internal/templates/pages/*.html",
        "internal/templates/components/*.html",
    }
    
    for _, pattern := range patterns {
        files, err := filepath.Glob(pattern)
        if err != nil {
            continue
        }
        if len(files) > 0 {
            tmpl, err = tmpl.ParseFiles(files...)
            if err != nil {
                fmt.Printf("ERROR parsing templates: %v\n", err)
            }
        }
    }
    
    fmt.Printf("DEBUG: Loaded templates: %v\n", tmpl.DefinedTemplates())
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")
    
    err := tmpl.ExecuteTemplate(w, name, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func SetupRoutes(mux *http.ServeMux, db *sql.DB) {
    authHandler := NewAuthHandler(db)
    userHandler := NewUserHandler(db)
    machineHandler := NewMachineHandler(db)
    dashboardHandler := NewDashboardHandler(db)
    locationHandler := NewLocationHandler(db)
    operationHandler := NewOperationHandler(db)

    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    
    // Auth routes
    mux.HandleFunc("/auth/signin", authHandler.SignIn)
    mux.HandleFunc("/auth/signup", authHandler.SignUp)
    mux.HandleFunc("/auth/signout", authHandler.SignOut)
    
    // Dashboard
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    })
    
    // Protected routes
    mux.HandleFunc("/dashboard", authHandler.RequireAuth(dashboardHandler.Dashboard))
    mux.HandleFunc("/users", authHandler.RequireAuth(userHandler.ListUsers))
    mux.HandleFunc("/users/form", authHandler.RequireAuth(userHandler.GetUserForm))
    mux.HandleFunc("/users/save", authHandler.RequireAuth(userHandler.SaveUser))
    mux.HandleFunc("/users/delete", authHandler.RequireAuth(userHandler.DeleteUser))
    mux.HandleFunc("/locations", authHandler.RequireAuth(locationHandler.ListLocations))
    mux.HandleFunc("/locations/form", authHandler.RequireAuth(locationHandler.GetLocationForm))
    mux.HandleFunc("/locations/save", authHandler.RequireAuth(locationHandler.SaveLocation))
    mux.HandleFunc("/locations/delete", authHandler.RequireAuth(locationHandler.DeleteLocation))
    mux.HandleFunc("/machines", authHandler.RequireAuth(machineHandler.ListMachines))
    mux.HandleFunc("/machines/form", authHandler.RequireAuth(machineHandler.GetMachineForm))
    mux.HandleFunc("/machines/save", authHandler.RequireAuth(machineHandler.SaveMachine))
    mux.HandleFunc("/machines/delete", authHandler.RequireAuth(machineHandler.DeleteMachine))
    mux.HandleFunc("/operations", authHandler.RequireAuth(operationHandler.ListOperations))
    mux.HandleFunc("/operations/form", authHandler.RequireAuth(operationHandler.GetOperationForm))
    mux.HandleFunc("/operations/save", authHandler.RequireAuth(operationHandler.SaveOperation))
    mux.HandleFunc("/api/stats", authHandler.RequireAuth(dashboardHandler.GetStats))
}
ENDOFFILE
fi

# Создаем минимальный main.go если нужно
if [ ! -f "main.go" ]; then
    echo "📝 Создаем main.go..."
    cat > main.go << 'ENDOFFILE'
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
ENDOFFILE
fi

# Проверяем наличие go.mod
if [ ! -f "go.mod" ]; then
    echo "📝 Инициализируем go.mod..."
    go mod init vend_erp
fi

# Добавляем необходимые зависимости
echo "📦 Добавляем зависимости..."
go get github.com/lib/pq
go mod tidy

# Проверяем компиляцию
echo "🔧 Проверяем компиляцию..."
if go build -o /tmp/vend_erp .; then
    echo "✅ Компиляция успешна!"
else
    echo "❌ Ошибка компиляции. Показываем детали:"
    go build -o /tmp/vend_erp . 2>&1
    exit 1
fi

echo "🚀 Запускаем приложение..."
air