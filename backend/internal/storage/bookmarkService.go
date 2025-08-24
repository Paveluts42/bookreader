package storage

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type BookmarkService struct {
    db *gorm.DB
}

func NewBookmarkService(db *gorm.DB) *BookmarkService {
    return &BookmarkService{db: db}
}

func (s *BookmarkService) AddBookmark(bookID, userID, note string, page int32) (*Bookmark, error) {
    bookmark := Bookmark{
        ID:     uuid.New(),
        BookID: uuid.MustParse(bookID),
        UserID: uuid.MustParse(userID),
        Page:   page,
        Note:   note,
    }
    if err := s.db.Create(&bookmark).Error; err != nil {
        return nil, err
    }
    return &bookmark, nil
}

func (s *BookmarkService) GetBookmarks(bookID, userID string) ([]Bookmark, error) {
    var bookmarks []Bookmark
    err := s.db.Where("book_id = ? AND user_id = ?", bookID, userID).Find(&bookmarks).Error
    return bookmarks, err
}

func (s *BookmarkService) DeleteBookmark(bookmarkID, userID string) error {
    return s.db.Where("id = ? AND user_id = ?", bookmarkID, userID).Delete(&Bookmark{}).Error
}