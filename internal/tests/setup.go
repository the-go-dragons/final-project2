package tests

import (
	"database/sql"
	"fmt"

	// _ "github.com/lib/pq"
	// "github.com/the-go-dragons/final-project2/internal/app"
	// "github.com/the-go-dragons/final-project2/pkg/database"

	"path"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/the-go-dragons/final-project2/internal/app"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"gorm.io/gorm"
)

var TestPgDb *sql.DB

var GormDb *gorm.DB

var RouteApp *app.App

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	config.Path = dir
}

func Setup() {
	var err error
	// config.LoadTestEnvVariables()
	config.Load()
	database.Load()
	// Open a connection to the test database
	testConStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Test,
		config.Config.Database.Ssl,
	)
	TestPgDb, err = sql.Open("postgres", testConStr)
	if err != nil {
		panic(err)
	}

	//Drop all tables
	_, err = TestPgDb.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		panic(err)
	}

	//Create databse connection and migrate tables
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
