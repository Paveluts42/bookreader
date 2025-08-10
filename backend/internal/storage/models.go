package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FilePath  string
	Title     string
	Author    string
	CoverPath string
	CreatedAt time.Time
    Page     int32     `json:"page"`     
    PageAll  int32     `json:"pageAll"`

	Notes     []Note     `gorm:"foreignKey:BookID"`
}

type Note struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookID uuid.UUID `gorm:"type:uuid"`
	Page   int
	Text   string
}



func AutoMigrate(DB *gorm.DB) {
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	DB.AutoMigrate(&Book{}, &Note{})
}
