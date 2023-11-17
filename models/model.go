package models

import (
	"database/sql"
	"fmt"

	"github.com/Thauan/golang-api-upload-s3-example/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Create an exported global variable to hold the database connection pool.
var DB *sql.DB

func dbConnection() (*gorm.DB, error) {
	handlers.LoadEnv()
	DatabasePort := handlers.GetEnvWithKey("DATABASE_PORT")
	DatabaseHost := handlers.GetEnvWithKey("DATABASE_HOST")
	DatabaseTable := handlers.GetEnvWithKey("DATABASE_NAME")
	DatabaseUser := handlers.GetEnvWithKey("DATABASE_USER")
	DatabasePassword := handlers.GetEnvWithKey("DATABASE_PASSWORD")
	sslMode := handlers.GetEnvWithKey("SSL_MODE")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		DatabaseHost, DatabaseUser, DatabasePassword, DatabaseTable, DatabasePort, sslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("could not insert row: %v", err)
		panic(err)
	}

	return db, err
}
