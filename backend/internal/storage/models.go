package storage

import _ "gorm.io/gorm"

type Book struct {
    ID       string `gorm:"primaryKey"`
    FilePath string
}

type Note struct {
    ID     string `gorm:"primaryKey"`
    BookID string
    Page   int
    Text   string
}

type Position struct {
    BookID      string `gorm:"primaryKey"`
    CurrentPage int
}

func AutoMigrate() {
    DB.AutoMigrate(&Book{}, &Note{}, &Position{})
}