package database

import (
	"base/src/business/entities"
	"fmt"

	"gorm.io/gorm"
)

func createTableMessages(db *gorm.DB) {
	err := db.AutoMigrate(entities.Messages{})
	if err != nil {
		fmt.Println(err)
		panic("error trying to migrate message schema")
	}

	fmt.Printf("Migration of Messages completed.\n\n")
}

func MigrateAllTables(db *gorm.DB) {
	createTableMessages(db)
}
