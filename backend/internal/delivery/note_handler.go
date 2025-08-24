package delivery

import (
    "context"
    "errors"
    "log"

    "connectrpc.com/connect"
    "github.com/Paveluts42/bookreader/backend/api"
    "github.com/Paveluts42/bookreader/backend/internal/shared"
    "github.com/Paveluts42/bookreader/backend/internal/storage"
    "gorm.io/gorm"
)

func (s *Server) AddNote(
    ctx context.Context,
    req *connect.Request[api.AddNoteRequest],
) (*connect.Response[api.AddNoteResponse], error) {
    log.Println("🔥 AddNote from CONNECT MAIN.go")
    userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }
    noteService := storage.NewNoteService(storage.DB)
    _, err = noteService.AddNote(req.Msg.BookId, userID, req.Msg.Text, req.Msg.Page)
    if err != nil {
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
    noteService := storage.NewNoteService(storage.DB)
    notes, err := noteService.GetNotes(req.Msg.BookId, userID)
    if err != nil {
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