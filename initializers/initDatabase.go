package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/prakyy/Go-Drive/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB;

func LinkDB() {
	var err error
	dsn := os.Getenv("DBConnectionString")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB")
	}
	fmt.Println("1/3 ✅ Database connection successful")
	//fmt.Print(DB)
}

func SyncDB() {
    if DB == nil {
        log.Fatal("Database connection is not initialized")
    }

    // Auto-migration for user model
    err := DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Auto migration failed user model", err)
    }
    fmt.Println("2/3 ✅ Database migration successful")

    // Auto migration for file model
    DB.AutoMigrate(&models.File{})
    if err != nil {
        log.Fatal("Auto migration failed file model:", err)
    }
    fmt.Println("3/3 ✅ File migration successful")
}