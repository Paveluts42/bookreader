package storage

import (
    "os"

    "gorm.io/gorm"
)

type BookService struct {
    db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
    return &BookService{db: db}
}

func (s *BookService) GetBooks(userID string, isAdmin bool) ([]Book, error) {
    var books []Book
    var err error
    if isAdmin {
        err = s.db.Find(&books).Error
    } else {
        err = s.db.Where("user_id = ?", userID).Find(&books).Error
    }
    return books, err
}

func (s *BookService) GetBook(bookID, userID string, isAdmin bool) (*Book, error) {
    var book Book
    var err error
    if isAdmin {
        err = s.db.First(&book, "id = ?", bookID).Error
    } else {
        err = s.db.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error
    }
    return &book, err
}

func (s *BookService) UpdateBookPage(bookID, userID string, isAdmin bool, page int32) error {
    var book Book
    var err error
    if isAdmin {
        err = s.db.First(&book, "id = ?", bookID).Error
    } else {
        err = s.db.Where("id = ? AND user_id = ?", bookID, userID).First(&book).Error
    }
    if err != nil {
        return err
    }
    book.Page = page
    return s.db.Save(&book).Error
}

func (s *BookService) DeleteBookWithData(bookID string) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("book_id = ?", bookID).Delete(&Note{}).Error; err != nil {
            return err
        }
        if err := tx.Where("book_id = ?", bookID).Delete(&Bookmark{}).Error; err != nil {
            return err
        }
        var book Book
        if err := tx.First(&book, "id = ?", bookID).Error; err != nil {
            return err
        }
        var user User
        if err := tx.First(&user, "id = ?", book.UserID).Error; err != nil {
            return err
        }
        bookDir := "/uploads/" + user.Username + "/" + bookID
        os.RemoveAll(bookDir)
        if err := tx.Delete(&Book{}, "id = ?", bookID).Error; err != nil {
            return err
        }

        return nil
    })
}