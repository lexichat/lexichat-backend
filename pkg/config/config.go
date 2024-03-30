package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
    DBHost     = os.Getenv("POSTGRES_HOST")
    DBPort     = os.Getenv("POSTGRES_PORT")
    DBUser     = os.Getenv("POSTGRES_USER")
    DBPassword = os.Getenv("POSTGRES_PASSWORD")
    DBName     = os.Getenv("POSTGRES_DB")
	GOLANG_SERVER_PORT = os.Getenv("GOLANG_SERVER_PORT")
    JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
)