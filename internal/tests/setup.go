package tests

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/database"

	"gorm.io/gorm"
)

var TestPgDb *sql.DB

var GormDb *gorm.DB

var RouteApp *app.App

func Setup() {
	var err error
	config.LoadTestEnvVariables()

	// Open a connection to the test database
	testConStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.GetEnv("POSTGRES_USER"),
		config.GetEnv("POSTGRES_PASSWORD"),
		config.GetEnv("DATABASE_HOST"),
		config.GetEnv("DATABASE_PORT"),
		config.GetEnv("POSTGRES_TEST_DB"),
		config.GetEnv("POSTGRES_SSL"),
	)
	TestPgDb, err = sql.Open("postgres", testConStr)
	if err != nil {
		panic(err)
	}

	// Drop all tables
	_, err = TestPgDb.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		panic(err)
	}

	// Create databse connection and migrate tables
	database.CreateTestDBConnection()
	database.AutoMigrateDB()
	GormDb, err = database.GetDatabaseConnection()
	if err != nil {
		panic(err)
	}

	// New webserver
	RouteApp = app.NewApp()
}

func CleanUp() {
	TestPgDb.Close()
	database.CloseDBConnection(GormDb)
}
