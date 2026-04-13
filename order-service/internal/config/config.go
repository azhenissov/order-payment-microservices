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
		DBHost:            getEnv("DB_HOST"),
		DBUser:            getEnv("DB_USER"),
		DBPassword:        getEnv("DB_PASSWORD"),
		DBName:            getEnv("DB_NAME"),
		DBPort:            getEnv("DB_PORT"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL"),
		ServerPort:        getEnv("SERVER_PORT"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("environment variable not set: " + key)
	}
	return value
}
