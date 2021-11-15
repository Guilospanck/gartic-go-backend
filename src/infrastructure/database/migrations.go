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

func createTableDrawers(db *gorm.DB) {
	err := db.AutoMigrate(entities.Drawers{})
	if err != nil {
		fmt.Println(err)
		panic("error trying to migrate drawers schema")
	}

	fmt.Println("Migration of Drawers completed.")
}

func MigrateAllTables(db *gorm.DB) {
	createTableMessages(db)
	createTableDrawers(db)
}
