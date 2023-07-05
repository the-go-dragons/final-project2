package main

import (
	"log"

	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

func main() {
	config.Load()
	database.Load()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	app := app.NewApp()
	// seeder.Run()
	rabbitmq.Connect()
	cronjob.NewCronJobRunner()
	log.Fatalln(app.Start(config.Config.Server.Port))
}
