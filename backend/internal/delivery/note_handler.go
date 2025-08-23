package delivery

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
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

    userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }

    var notes []storage.Note
    if err := storage.DB.Where("book_id = ? AND user_id = ?", req.Msg.BookId, userID).Find(&notes).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, connect.NewError(connect.CodeNotFound, err)
        }
        log.Printf("DB error: %v", err)
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    apiNotes := make([]*api.Note, len(notes))
    for i, n := range notes {
        apiNotes[i] = &api.Note{
            Id:     n.ID.String(),
            Text:   n.Text,
            Page:   int32(n.Page),
            BookId: n.BookID.String(),
            UserId: n.UserID.String(),
        }
    }

    resp := &api.GetNotesResponse{Notes: apiNotes}
    return connect.NewResponse(resp), nil
}

