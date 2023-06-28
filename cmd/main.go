package main

import (
	"log"

	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

func main() {
	config.Load()
	database.Load()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	app := app.NewApp()
	// seeder.Run()
	log.Fatalln(app.Start(config.Config.Server.Port))
}
