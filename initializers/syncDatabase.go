package initializers

import (
	"log"

	"github.com/prakyy/Go-Drive/models"
)

func SyncDatabase() {
    if DB == nil {
        log.Fatal("Database connection is not initialized")
    }

    // Auto-migration for user model
    err := DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Auto migration failed user model", err)
    }
    log.Println("Database migration successful")

    // Auto migration for file model
    DB.AutoMigrate(&models.File{})
    if err != nil {
        log.Fatal("Auto migration failed file model:", err)
    }
    log.Println("File model migration successful")
}