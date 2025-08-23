package storage

import (
	"log"
	"os"

	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	DB = db
	log.Println("Connected to PostgreSQL via GORM")
	AutoMigrate(DB)
    createAdminUser()
}

func createAdminUser() {
	var count int64
	DB.Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		passwordHash, err := shared.HashPassword("admin")
		if err != nil {
			log.Fatalf("failed to hash admin password: %v", err)
		}

		admin := User{
			Username:     "admin",
			PasswordHash: passwordHash,
			IsAdmin:      true,
		}
		DB.Create(&admin)
	}
}
