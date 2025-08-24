package storage

import (
	"os"
	"time"

	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func (s *UploadService) SavePDF(bookID, title, author, userID string, chunk []byte, pageCount int, filePath, coverPath string) (*Book, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if _, err := file.Write(chunk); err != nil {
		return nil, err
	}
	if err := shared.GenerateCover(filePath, coverPath); err != nil {
		return nil, err
	}
	book := Book{
		ID:        uuid.MustParse(bookID),
		Title:     title,
		Author:    author,
		Page:      int32(0),
		PageAll:   int32(pageCount),
		FilePath:  filePath,
		CoverPath: coverPath,
		UserID:    uuid.MustParse(userID),
		CreatedAt: time.Now(),
	}
	if err := DB.Create(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}
