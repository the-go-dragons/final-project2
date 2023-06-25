package main

import (
	"fmt"
	"os"

	"github.com/the-go-dragons/final-project2/pkg/database"
)

func main() {
	// config.LoadEnvVariables()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
}
