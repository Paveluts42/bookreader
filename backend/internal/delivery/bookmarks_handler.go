package delivery

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
)




func (s *Server) AddBookmark(
	ctx context.Context,
	req *connect.Request[api.AddBookmarkRequest],
) (*connect.Response[api.AddBookmarkResponse], error) {
    userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }
    bookmarkService := storage.NewBookmarkService(storage.DB)
    _, err = bookmarkService.AddBookmark(req.Msg.BookId, userID, req.Msg.Note, req.Msg.Page)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }
    resp := &api.AddBookmarkResponse{Ok: true}
    return connect.NewResponse(resp), nil
}


func (s *Server) GetBookmarks(
	ctx context.Context,
	req *connect.Request[api.GetBookmarksRequest],
) (*connect.Response[api.GetBookmarksResponse], error) {
 userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }
    bookmarkService := storage.NewBookmarkService(storage.DB)
    bookmarks, err := bookmarkService.GetBookmarks(req.Msg.BookId, userID)
    if err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }
    apiBookmarks := make([]*api.Bookmark, len(bookmarks))
    for i, b := range bookmarks {
        apiBookmarks[i] = &api.Bookmark{
            Id:     b.ID.String(),
            BookId: b.BookID.String(),
            Page:   int32(b.Page),
            Note:   b.Note,
            UserId: b.UserID.String(),
        }
    }
    resp := &api.GetBookmarksResponse{Bookmarks: apiBookmarks}
    return connect.NewResponse(resp), nil
}

func (s *Server) DeleteBookmark(
	ctx context.Context,
	req *connect.Request[api.DeleteBookmarkRequest],
) (*connect.Response[api.DeleteBookmarkResponse], error) {
    userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }
    bookmarkService := storage.NewBookmarkService(storage.DB)
    if err := bookmarkService.DeleteBookmark(req.Msg.BookmarkId, userID); err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }
    resp := &api.DeleteBookmarkResponse{Ok: true}
    return connect.NewResponse(resp), nil
}
