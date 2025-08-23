package delivery

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"github.com/google/uuid"
)

func (s *Server) UploadPDF(
	ctx context.Context,
	req *connect.Request[api.UploadRequest],
) (*connect.Response[api.UploadResponse], error) {

    userID, err := shared.ValidateAccessToken(req)
	println("UserID:", userID)
    if err != nil || userID == "" {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }

	if req.Msg.BookId == "" || len(req.Msg.Chunk) == 0 {
		log.Println("Invalid upload request")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid upload request"))
	}

	title := req.Msg.Title
	author := req.Msg.Author

	bookID := req.Msg.BookId

	filePath := "/uploads/" + bookID + ".pdf"
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer file.Close()

	if _, err := file.Write(req.Msg.Chunk); err != nil {
		log.Printf("Failed to write file: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	coverPath := "/uploads/" + bookID + ".png"
	if err := shared.GenerateCover(filePath, coverPath); err != nil {
		log.Printf("Failed to generate cover: %v", err)
	}

	pageCount, err := shared.GetPDFPageCount(filePath)
	if err != nil {
		log.Printf("Failed to get page count: %v", err)
		pageCount = 0
	}
	book := storage.Book{
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
	if err := storage.DB.Create(&book).Error; err != nil {
		log.Printf("Failed to save book: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.UploadResponse{BookId: bookID, Ok: true}
	return connect.NewResponse(resp), nil
}
