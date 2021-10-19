package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_DATABASE_NAME")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, username, password, databaseName, dbPort, sslMode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	fmt.Printf("%s", fmt.Sprintf("Connected to the database %s at %s:%s!\n", databaseName, host, dbPort))

	/* Migrate tables */
	MigrateAllTables(db)

	return db
}
