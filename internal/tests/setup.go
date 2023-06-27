package tests

import (
	_ "github.com/lib/pq"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"github.com/the-go-dragons/final-project2/pkg/database"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Setup() {
	var err error
	config.LoadEnvVariables()
	database.CreateDBConnection()
	database.AutoMigrateDB()
	Db, err = database.GetDatabaseConnection()
	_ = err
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
