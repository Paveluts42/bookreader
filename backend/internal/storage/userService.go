package storage

import (
    "log"
    "os"

    "gorm.io/gorm"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}


func (s *UserService) IsAdmin(userID string) bool {
    var user User
    if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
        return false
    }
    return user.IsAdmin
}
func (s *UserService) GetUsers() ([]User, error) {
    var users []User
    err := s.db.Find(&users).Error
    return users, err
}

func (s *UserService) DeleteUserWithData(userID string) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        var books []Book
        if err := tx.Where("user_id = ?", userID).Find(&books).Error; err != nil {
            return err
        }
        for _, book := range books {
            if err := tx.Where("book_id = ?", book.ID).Delete(&Note{}).Error; err != nil {
                return err
            }
            if err := tx.Where("book_id = ?", book.ID).Delete(&Bookmark{}).Error; err != nil {
                return err
            }
            if err := tx.Delete(&book).Error; err != nil {
                return err
            }
            if book.FilePath != "" {
                if err := os.Remove(book.FilePath); err != nil && !os.IsNotExist(err) {
                    log.Printf("Failed to delete PDF file %s: %v", book.FilePath, err)
                }
            }
            if book.CoverPath != "" {
                if err := os.Remove(book.CoverPath); err != nil && !os.IsNotExist(err) {
                    log.Printf("Failed to delete PNG file %s: %v", book.CoverPath, err)
                }
            }
        }
        if err := tx.Delete(&User{}, "id = ?", userID).Error; err != nil {
            return err
        }
        return nil
    })
}

