package main

import (
	"authenservice/database"
	"authenservice/database/models"
	"log"
)

// RunMigrations ทำการ migrate ตาราง
func main() {
	database.Connect()

	// ทำการ AutoMigrate
	err := database.DB.AutoMigrate(
		&models.Role{},
		&models.UserAuth{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("Migration completed successfully")
}
