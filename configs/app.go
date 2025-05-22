package configs

import (
	"restful-api/pkg/utls"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppHost string
	AppPort int
}

func NewAppConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("API_HOST")
	port := utls.GetenvInt("API_PORT", 8585)
	return &AppConfig{
		AppHost: host,
		AppPort: port,
	}
}
