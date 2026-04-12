package config

import (
	"os"
)

type Config struct {
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBPort            string
	PaymentServiceURL string
	ServerPort        string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBUser:            getEnv("DB_USER", "order_user"),
		DBPassword:        getEnv("DB_PASSWORD", "order_password"),
		DBName:            getEnv("DB_NAME", "orders_db"),
		DBPort:            getEnv("DB_PORT", "5434"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8081"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
