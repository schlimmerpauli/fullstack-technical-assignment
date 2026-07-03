package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// this usualy contains the db connection strings and other configuration values that are needed
type Config struct {
	ServerHost         string
	ServerPort         string
	CORSAllowedOrigins []string
}

const (
	DefaultServerHost         = "0.0.0.0"
	DefaultServerPort         = "8080"
	DefaultCORSAllowedOrigins = "*"
)

func LoadConfig() Config {
	loadDotEnv()
	cfg := Config{
		ServerHost:         DefaultServerHost,
		ServerPort:         DefaultServerPort,
		CORSAllowedOrigins: splitCSV(DefaultCORSAllowedOrigins),
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.ServerHost = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	}

	if origins := os.Getenv("CORS_ALLOWED_ORIGINS"); origins != "" {
		cfg.CORSAllowedOrigins = splitCSV(origins)
	}

	return cfg
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		items = append(items, item)
	}

	return items
}

func loadDotEnv() {
	path := ".env"
	if err := godotenv.Load(path); err != nil {
		log.Printf("No .env file found at %s, using default values", path)
	}
}
