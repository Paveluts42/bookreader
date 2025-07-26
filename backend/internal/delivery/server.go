package delivery

import (
	"context"

	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
)

type server struct {
	api.UnimplementedReaderServiceServer
}

func (s *server) SaveNote(ctx context.Context, req *api.Note) (*api.NoteResponse, error) {
	note := storage.Note{
		ID:     req.Id,
		BookID: req.BookId,
		Page:   int(req.Page),
		Text:   req.Text,
	}
	if err := storage.DB.Save(&note).Error; err != nil {
		return &api.NoteResponse{Ok: false}, err
	}
	return &api.NoteResponse{Ok: true, NoteId: note.ID}, nil
}

func (s *server) ListNotes(req *api.ListNotesRequest, stream api.ReaderService_ListNotesServer) error {
	var notes []storage.Note
	if err := storage.DB.Where("book_id = ?", req.BookId).Find(&notes).Error; err != nil {
		return err
	}
	for _, n := range notes {
		stream.Send(&api.Note{
			Id:     n.ID,
			BookId: n.BookID,
			Page:   int32(n.Page),
			Text:   n.Text,
		})
	}
	return nil
}

func (s *server) SavePosition(ctx context.Context, req *api.ReadingPosition) (*api.PositionResponse, error) {
	pos := storage.Position{
		BookID:      req.BookId,
		CurrentPage: int(req.CurrentPage),
	}
	if err := storage.DB.Save(&pos).Error; err != nil {
		return &api.PositionResponse{Ok: false}, err
	}
	return &api.PositionResponse{Ok: true}, nil
}

func (s *server) GetPosition(ctx context.Context, req *api.ListNotesRequest) (*api.ReadingPosition, error) {
	var pos storage.Position
	if err := storage.DB.First(&pos, "book_id = ?", req.BookId).Error; err != nil {
		return &api.ReadingPosition{BookId: req.BookId, CurrentPage: 0}, nil
	}
	return &api.ReadingPosition{BookId: pos.BookID, CurrentPage: int32(pos.CurrentPage)}, nil
}
