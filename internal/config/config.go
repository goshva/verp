package config

import (
    "os"
    "strings"
)

type Config struct {
    DBHost        string
    DBPort        string
    DBUser        string
    DBPassword    string
    DBName        string
    ServerAddress string
    SessionSecret string
    AppEnv        string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:        getEnv("DB_HOST", "master.28100775-6c4d-4114-b141-1399e7cbef21.c.dbaas.selcloud.ru"),
        DBPort:        getEnv("DB_PORT", "5432"),
        DBUser:        getEnv("DB_USER", "djqme"),
        DBPassword:    getEnv("DB_PASSWORD", "dFHdsDhwZXeN"),
        DBName:        getEnv("DB_NAME", "12xxxdy6c"),
        ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
        SessionSecret: getEnv("SESSION_SECRET", "your-secret-key-change-in-production"),
        AppEnv:        getEnv("APP_ENV", "development"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

// LoadEnv loads environment variables from .env file (for development)
func LoadEnv() error {
    if _, err := os.Stat(".env"); err == nil {
        content, err := os.ReadFile(".env")
        if err != nil {
            return err
        }
        
        lines := strings.Split(string(content), "\n")
        for _, line := range lines {
            line = strings.TrimSpace(line)
            if line == "" || strings.HasPrefix(line, "#") {
                continue
            }
            
            parts := strings.SplitN(line, "=", 2)
            if len(parts) == 2 {
                key := strings.TrimSpace(parts[0])
                value := strings.TrimSpace(parts[1])
                os.Setenv(key, value)
            }
        }
    }
    return nil
}