package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// this usualy contains the db connection strings and other configuration values that are needed
type Config struct {
	ServerHost string
	ServerPort string
}

const (
	DefaultServerHost = "0.0.0.0"
	DefaultServerPort = "8080"
)

func LoadConfig() Config {
	loadDotEnv()
	cfg := Config{
		ServerHost: DefaultServerHost,
		ServerPort: DefaultServerPort,
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		cfg.ServerHost = host
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.ServerPort = port
	}
	return cfg
}

func loadDotEnv() {
	path := ".env"
	if err := godotenv.Load(path); err != nil {
		log.Printf("No .env file found at %s, using default values", path)
	}
}
