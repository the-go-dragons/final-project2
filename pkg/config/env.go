package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func LoadTestEnvVariables() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
}

// This function can be used to get ENV Var with default value
func GetEnv(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatal(key + " is not found from env files!")
		return ""
	}
	return value
}
