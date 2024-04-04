package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	APP_SERVER_PORT  string
	WS_SERVER_PORT   string
	JWT_SECRET_KEY   string
)

func LoadEnvVariables(envPath string) {
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	DBHost = os.Getenv("POSTGRES_HOST")
	DBPort = os.Getenv("POSTGRES_PORT")
	DBUser = os.Getenv("POSTGRES_USER")
	DBPassword = os.Getenv("POSTGRES_PASSWORD")
	DBName = os.Getenv("POSTGRES_DB")
	APP_SERVER_PORT = os.Getenv("APP_SERVER_PORT")
	WS_SERVER_PORT = os.Getenv("WS_SERVER_PORT")
	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
}