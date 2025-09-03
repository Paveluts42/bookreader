package delivery

import (
	"context"
	"errors"
	"log"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
)

func (s *Server) UploadPDF(
    ctx context.Context,
    req *connect.Request[api.UploadRequest],
) (*connect.Response[api.UploadResponse], error) {
    userID, err := shared.ValidateAccessToken(req)
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

    pageCount, err := shared.GetPDFPageCountFromBytes(req.Msg.Chunk)
    if err != nil {
        log.Printf("Failed to get page count: %v", err)
        pageCount = 0
    }

    uploadService := storage.NewUploadService()
    _, err = uploadService.SavePDF(bookID, title, author, userID, req.Msg.Chunk, pageCount, "", "")
    if err != nil {
        log.Printf("Failed to save book: %v", err)
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    resp := &api.UploadResponse{BookId: bookID, Ok: true}
    return connect.NewResponse(resp), nil
}