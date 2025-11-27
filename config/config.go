package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
    _ "github.com/jackc/pgx/v4/stdlib" 
)

type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    SSLMode    string
}

func LoadConfig() *Config {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Printf("Warning: Error loading .env file: %v", err)
        // Continue with environment variables or defaults
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnvAsInt("DB_PORT", 5432),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "venderp"),
        SSLMode:    getEnv("SSL_MODE", "disable"),
    }
    
    return config
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func (c *Config) GetConnectionString() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)
}

func ConnectDB(config *Config) (*sql.DB, error) {
    // Remove problematic environment variable
    os.Unsetenv("PGLOCALEDIR")
    
    connStr := config.GetConnectionString()
    
    // Use "pgx" as driver name instead of "postgres"
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %v", err)
    }

    // Test the connection
    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("error connecting to database: %v", err)
    }

    log.Printf("Successfully connected to database: %s", config.DBName)
    return db, nil
}