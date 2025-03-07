package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит настройки приложения.
type Config struct {
	GRPCPort int
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string
}

// LoadConfig загружает конфигурацию из `.env` файла и переменных окружения.
func LoadConfig() (*Config, error) {
	// Загружаем переменные окружения из .env файла (если есть)
	if err := godotenv.Load(); err != nil {
		fmt.Println("Нет .env файла, используем переменные окружения")
	}

	grpcPort, err := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации GRPC_PORT: %v", err)
	}

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации DB_PORT: %v", err)
	}

	cfg := &Config{
		GRPCPort: grpcPort,
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   dbPort,
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASS", "password"),
		DBName:   getEnv("DB_NAME", "testdb"),
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или дефолтное значение, если переменная не установлена.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
