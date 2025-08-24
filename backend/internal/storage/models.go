package storage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	IsAdmin   bool      `gorm:"default:false"`
	Books     []Book     `gorm:"foreignKey:UserID"`
}

type Book struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FilePath  string
	Title     string
	Author    string
	CoverPath string
	CreatedAt time.Time
    Page     int32     `json:"page"`     
    PageAll  int32     `json:"pageAll"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Notes     []Note     `gorm:"foreignKey:BookID"`
}

type Note struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookID uuid.UUID `gorm:"type:uuid"`
	Page   int
	Text   string
	UserID uuid.UUID `gorm:"type:uuid"`
}
type Bookmark struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookID uuid.UUID `gorm:"type:uuid"`
	UserID uuid.UUID `gorm:"type:uuid"`
	Page   int32
	Note   string
}

func AutoMigrate(DB *gorm.DB) {
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	DB.AutoMigrate(&User{}, &Book{}, &Note{}, &Bookmark{})
}
