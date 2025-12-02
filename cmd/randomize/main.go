// cmd/randomize/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Config структура конфигурации из переменных окружения
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() (*Config, error) {
	// Пытаемся загрузить .env файл, но не падаем если его нет
	_ = godotenv.Load(".env")
	
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		port = 5432
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     port,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "vend_erp"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
	}

	// Проверяем обязательные поля
	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD не установлен в .env файле")
	}

	return cfg, nil
}

// getEnv вспомогательная функция для получения переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// ConnectDB устанавливает соединение с базой данных
func ConnectDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

// generateRandomTime генерирует случайное время в пределах последнего месяца
func generateRandomTime() time.Time {
	now := time.Now()
	// Начало месяца (30 дней назад)
	startOfMonth := now.AddDate(0, 0, -30)
	
	// Случайная разница в секундах
	diffSeconds := int64(now.Sub(startOfMonth).Seconds())
	randomSeconds := rand.Int63n(diffSeconds)
	
	return startOfMonth.Add(time.Duration(randomSeconds) * time.Second)
}

// updateTableDates обновляет даты в конкретной таблице
func updateTableDates(db *sql.DB, tableName string, hasUpdatedAt bool) error {
	// Получаем все ID из таблицы
	rows, err := db.Query(fmt.Sprintf("SELECT id FROM %s ORDER BY id", tableName))
	if err != nil {
		return fmt.Errorf("error querying %s: %w", tableName, err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error scanning id: %w", err)
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		log.Printf("Таблица %s пуста, пропускаем", tableName)
		return nil
	}

	log.Printf("Найдено %d записей в таблице %s", len(ids), tableName)

	// Обновляем каждую запись с транзакцией для больших таблиц
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	updatedCount := 0
	for _, id := range ids {
		createdAt := generateRandomTime()
		updatedAt := createdAt
		
		// Если есть updated_at и с 30% вероятностью, делаем его позже
		if hasUpdatedAt && rand.Intn(100) < 30 {
			updatedAt = createdAt.Add(time.Duration(rand.Intn(86400)) * time.Second) // до 24 часов позже
		}

		var query string
		var args []interface{}
		
		if hasUpdatedAt {
			query = fmt.Sprintf("UPDATE %s SET created_at = $1, updated_at = $2 WHERE id = $3", tableName)
			args = []interface{}{createdAt, updatedAt, id}
		} else {
			query = fmt.Sprintf("UPDATE %s SET created_at = $1 WHERE id = $2", tableName)
			args = []interface{}{createdAt, id}
		}

		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("error updating %s id=%d: %w", tableName, id, err)
		}
		updatedCount++

		// Логируем прогресс каждые 100 записей
		if updatedCount%100 == 0 {
			log.Printf("  Прогресс: %d/%d обновлено", updatedCount, len(ids))
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction for %s: %w", tableName, err)
	}

	log.Printf("✅ Обновлено %d записей в таблице %s", updatedCount, tableName)
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	log.Println("🚀 Запуск скрипта рандомизации дат...")
	
	// Загружаем конфигурацию из .env
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}

	// Подключаемся к базе данных
	db, err := ConnectDB(cfg)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	log.Println("✅ Подключение к базе данных успешно")
	log.Printf("📊 База данных: %s@%s:%d/%s", 
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Список таблиц и их структура (есть ли updated_at)
	tables := []struct {
		name          string
		hasUpdatedAt  bool
	}{
		// Таблицы с обоими полями
		{"public.locations", true},
		{"public.users", true},
		{"public.vending_machines", true},
		{"public.vending_operations", true},
		{"public.warehouse", true},
		{"public.warehouse_inventory", true},
		{"public.warehouse_shipments", true},
		{"public.warehouse_supplies", true},
		
		// Таблицы только с created_at
		{"public.inventory_adjustments", false},
		{"public.inventory_transfers", false},
		{"public.schema_migrations", false},
		{"public.sessions", false},
		{"public.shipment_items", false},
		{"public.supply_items", false},
		{"public.warehouse_categories", false},
	}

	// Подтверждение пользователя
	log.Println("⚠️  ВНИМАНИЕ: Этот скрипт изменит даты в базе данных!")
	log.Println("   Нажмите Enter для продолжения или Ctrl+C для отмены...")
	fmt.Scanln()

	// Обновляем каждую таблицу
	startTime := time.Now()
	for _, table := range tables {
		log.Printf("\n🔄 Обновление таблицы %s...", table.name)
		if err := updateTableDates(db, table.name, table.hasUpdatedAt); err != nil {
			log.Printf("⚠️  Ошибка при обновлении %s: %v", table.name, err)
		}
		time.Sleep(100 * time.Millisecond) // Небольшая пауза между таблицами
	}

	elapsed := time.Since(startTime)
	log.Printf("\n🎉 Все даты успешно рандомизированы за %v!", elapsed)
	log.Println("✅ Готово!")
}