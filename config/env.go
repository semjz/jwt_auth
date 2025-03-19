package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnvVariable loads a value from the .env file
func LoadEnvVariable(key string) string {
	// Load .env file
	err := godotenv.Load("D:/golang_programs/jwt_auth/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return the environment variable value
	return os.Getenv(key)
}
