/*	This module has two kind of databse connection.
	One is for main project and the other is for testing.
*/

package database

import (
	"fmt"
	"log"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	model "github.com/the-go-dragons/final-project2/internal/domain"
	"github.com/the-go-dragons/final-project2/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB

	db       string
	host     string
	port     int
	ssl      string
	timezone string
	user     string
	password string

	testDb string
)

func Load() {
	user = config.Config.Database.User
	password = config.Config.Database.Password
	db = config.Config.Database.Name
	host = config.Config.Database.Host
	port = config.Config.Database.Port
	ssl = config.Config.Database.Ssl
	timezone = config.Config.Database.Timezone
	testDb = config.Config.Database.Test
	// user = "root"
	// password = "PSUliSks8J3QPrDuGIx9egwo"
	// db = "smsproject"
	// host = "luca.iran.liara.ir"
	// port = 31835
	// ssl = "disable"
	// timezone = "Asia/tehran"
	// testDb = "smsprojecttest"
}

func GetDSN() string {
	conStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, db, port, ssl, timezone)
	return conStr
}

func GetTestDSN() string {
	conStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, testDb, port, ssl, timezone)
	return conStr
}

func CreateDBConnection() error {
	// Close the existing connection if open
	if dbConn != nil {
		CloseDBConnection(dbConn)
	}

	db_conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db_conn.DB()

	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db_conn
	return err
}

func CreateTestDBConnection() error {
	// Close the existing connection if open
	if dbConn != nil {
		CloseDBConnection(dbConn)
	}
	println("dns:", GetTestDSN())
	db_conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  GetTestDSN(),
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db_conn.DB()

	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
	dbConn = db_conn
	return err
}

func GetDatabaseConnection() (*gorm.DB, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		return dbConn, err
	}
	if err := sqlDB.Ping(); err != nil {
		return dbConn, err
	}
	return dbConn, nil
}

func CloseDBConnection(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
}

func AutoMigrateDB() error {
	conn, err := GetDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = conn.AutoMigrate(
		&model.Number{},
		&model.User{},
		&model.Wallet{},
		&model.Payment{},
		&model.Transaction{},
		&model.Subscription{},
		&model.SMSTemplate{},
		&model.PhoneBook{},
		&model.SMSHistory{},
		&model.Contact{},
	)

	// sqlDB, err := conn.DB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// driver, err := pg.WithInstance(sqlDB, &pg.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// migrate, err := migrate.NewWithDatabaseInstance(
	// 	"file://./pkg/database/migrations",
	// 	"postgres", driver)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// migrate.Up()
	return err
}
