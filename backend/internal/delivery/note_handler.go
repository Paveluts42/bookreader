package delivery

import (
	"context"
	"errors"
	"log"


	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)




func (s *Server) AddNote(
	ctx context.Context,
	req *connect.Request[api.AddNoteRequest],
) (*connect.Response[api.AddNoteResponse], error) {
	log.Println("🔥 AddNote from CONNECT MAIN.go")
	note := storage.Note{
		ID:     uuid.New(),
		BookID: uuid.MustParse(req.Msg.BookId),
		Page:   int(req.Msg.Page),
		Text:   req.Msg.Text,
	}

	if err := storage.DB.Create(&note).Error; err != nil {
		log.Printf("DB error: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.AddNoteResponse{Ok: true}
	return connect.NewResponse(resp), nil
}


func (s *Server) GetNotes(
	ctx context.Context,
	req *connect.Request[api.GetNotesRequest],
) (*connect.Response[api.GetNotesResponse], error) {
	log.Println("🔥 GetNotes from CONNECT MAIN.go")
	var book storage.Book
	if err := storage.DB.Preload("Notes").First(&book, "id = ?", req.Msg.BookId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		log.Printf("DB error: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	apiNotes := make([]*api.Note, len(book.Notes))
	for i, n := range book.Notes {
		apiNotes[i] = &api.Note{
			Id:     n.ID.String(),
			Text:   n.Text,
			Page:   int32(n.Page),
			BookId: n.BookID.String(),
		}
	}

	resp := &api.GetNotesResponse{Notes: apiNotes}
	return connect.NewResponse(resp), nil
}


