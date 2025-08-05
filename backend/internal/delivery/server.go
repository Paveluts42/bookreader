package delivery

import (
	"context"
	"fmt"
	"os"

	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"github.com/google/uuid"
)

type server struct {
	api.UnimplementedReaderServiceServer
}

func NewServer() *server {
	return &server{}
}

func ParseUuid(id string) (uuid.UUID, error) {
	value, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return value, nil
}

func (s *server) SaveNote(ctx context.Context, req *api.Note) (*api.NoteResponse, error) {

	bookID, _ := ParseUuid(req.BookId)

	ID, _ := ParseUuid(req.Id)

	note := storage.Note{
		ID:     ID,
		BookID: bookID,
		Page:   int(req.Page),
		Text:   req.Text,
	}
	if err := storage.DB.Save(&note).Error; err != nil {
		return &api.NoteResponse{Ok: false}, err
	}
	return &api.NoteResponse{Ok: true, NoteId: note.ID.String()}, nil
}

func (s *server) UploadPDF(ctx context.Context, req *api.UploadRequest) (*api.UploadResponse, error) {
	// Парсим UUID
	bookID, err := uuid.Parse(req.BookId)
	if err != nil {
		return nil, err
	}

	// Путь к файлу
	filePath := fmt.Sprintf("../uploads/%s.pdf", bookID.String())

	// Открываем файл на дозапись (создаст, если нет)
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Пишем чанк
	_, err = f.Write(req.Chunk)
	if err != nil {
		return nil, err
	}

	return &api.UploadResponse{
		BookId: req.BookId,
		Ok:     true,
	}, nil
}

func (s *server) ListNotes(req *api.ListNotesRequest, stream api.ReaderService_ListNotesServer) error {
	var notes []storage.Note
	if err := storage.DB.Where("book_id = ?", req.BookId).Find(&notes).Error; err != nil {
		return err
	}
	for _, n := range notes {
		stream.Send(&api.Note{
			Id:     n.ID.String(),
			BookId: n.BookID.String(),
			Page:   int32(n.Page),
			Text:   n.Text,
		})
	}
	return nil
}

func (s *server) SavePosition(ctx context.Context, req *api.ReadingPosition) (*api.PositionResponse, error) {
	bookID, _ := ParseUuid(req.BookId)

	pos := storage.Position{
		BookID:      bookID,
		CurrentPage: int(req.CurrentPage),
	}
	if err := storage.DB.Save(&pos).Error; err != nil {
		return &api.PositionResponse{Ok: false}, err
	}
	return &api.PositionResponse{Ok: true}, nil
}

func (s *server) GetPosition(ctx context.Context, req *api.GetPDFRequest) (*api.ReadingPosition, error) {
	var pos storage.Position
	if err := storage.DB.First(&pos, "book_id = ?", req.BookId).Error; err != nil {
		return &api.ReadingPosition{BookId: req.BookId, CurrentPage: 0}, nil
	}
	return &api.ReadingPosition{BookId: pos.BookID.String(), CurrentPage: int32(pos.CurrentPage)}, nil
}

func (s *server) GetBooks(ctx context.Context, req *api.GetBooksRequest) (*api.GetBooksResponse, error) {
	var books []storage.Book
	if err := storage.DB.Find(&books).Error; err != nil {
		return nil, err
	}
	resp := &api.GetBooksResponse{}
	for _, b := range books {
		resp.Books = append(resp.Books, &api.Book{
			Id:     b.ID.String(),
			Title:  b.Title,
			Author: b.Author,
		})
	}
	return resp, nil
}
