package main

import (
	"ahsfnu-media-cloud/internal/database"
	"ahsfnu-media-cloud/internal/models"
	"log"
)

func main() {
	database.Init()

	db := database.GetDB()
	err := db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
