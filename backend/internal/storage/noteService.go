package storage

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type NoteService struct {
    db *gorm.DB
}

func NewNoteService(db *gorm.DB) *NoteService {
    return &NoteService{db: db}
}

func (s *NoteService) AddNote(bookID, userID, text string, page int32) (*Note, error) {
    note := Note{
        ID:     uuid.New(),
        BookID: uuid.MustParse(bookID),
        Page:   int(page),
        Text:   text,
        UserID: uuid.MustParse(userID),
    }
    if err := s.db.Create(&note).Error; err != nil {
        return nil, err
    }
    return &note, nil
}

func (s *NoteService) GetNotes(bookID, userID string) ([]Note, error) {
    var notes []Note
    err := s.db.Where("book_id = ? AND user_id = ?", bookID, userID).Find(&notes).Error
    return notes, err
}