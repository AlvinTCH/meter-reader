package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GlobDB *gorm.DB
var GlobSqlDB *sql.DB

func Migrate() {
	GlobDB.AutoMigrate(&MeterReadings{})
}

func Open() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), 5432, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	// with batch size for batched create
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})

	GlobDB = DB

	if err != nil {
		log.Fatal("Error connecting to database", err)
		return err
	}

	fmt.Println("Connected to database", DB)

	// Get generic database object sql.DB to use its functions
	SqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Error getting sql.DB object", err)
		return err
	}

	GlobSqlDB = SqlDB

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	SqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	SqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	SqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
