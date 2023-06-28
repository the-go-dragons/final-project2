package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"github.com/the-go-dragons/final-project2/pkg/rabbitmq"
)

func main() {
	config.LoadEnvVariables()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	app := app.NewApp()
	// seeder.Run()
	cronjob.NewCronJobRunner()
	rabbitmq.Connect()
	log.Info(app.Start(config.GetEnv("EXPOSE_PORT")))
}
