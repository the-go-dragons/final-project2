package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/cronjob"
	"github.com/the-go-dragons/final-project2/pkg/database"
)

func main() {
	config.LoadEnvVariables()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	app := app.NewApp()
	// seeder.Run()
	cronjob.NewCronJobRunner()
	log.Info(app.Start(config.GetEnv("EXPOSE_PORT")))
}
