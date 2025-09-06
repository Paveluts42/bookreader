package delivery

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"gorm.io/gorm"
)

func (s *Server) GetBooks(
	ctx context.Context,
	req *connect.Request[api.GetBooksRequest],
) (*connect.Response[api.GetBooksResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil || userID == "" {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
	}

	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	bookService := storage.NewBookService(storage.DB)
	books, err := bookService.GetBooks(userID, user.IsAdmin)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.GetBooksResponse{Books: make([]*api.Book, 0, len(books))}
	for _, b := range books {
		resp.Books = append(resp.Books, &api.Book{
			Id:        b.ID.String(),
			Title:     b.Title,
			Author:    b.Author,
			FilePath:  b.FilePath,
			CoverUrl:  b.CoverPath,
			AudioPath: b.AudioPath,
			UserId:    b.UserID.String(),
			CreatedAt: b.CreatedAt.Format("2006-01-02 15:04:05"),
			Page:      int32(b.Page),
			PageAll:   int32(b.PageAll),
		})
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) GetBook(
	ctx context.Context,
	req *connect.Request[api.GetBookRequest],
) (*connect.Response[api.GetBookResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
	}
	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	bookService := storage.NewBookService(storage.DB)
	book, err := bookService.GetBook(req.Msg.BookId, userID, user.IsAdmin)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	bookResp := &api.Book{
		Id:        book.ID.String(),
		Title:     book.Title,
		Author:    book.Author,
		FilePath:  book.FilePath,
		CoverUrl:  book.CoverPath,
		AudioPath: book.AudioPath,
		CreatedAt: book.CreatedAt.Format("2006-01-02 15:04:05"),
		Page:      int32(book.Page),
		PageAll:   int32(book.PageAll),
	}
	resp := &api.GetBookResponse{Book: bookResp}
	return connect.NewResponse(resp), nil
}

func (s *Server) DeleteBook(
	ctx context.Context,
	req *connect.Request[api.DeleteBookRequest],
) (*connect.Response[api.DeleteBookResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
	}
	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	bookService := storage.NewBookService(storage.DB)
	book, err := bookService.GetBook(req.Msg.BookId, userID, user.IsAdmin)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("book not found"))
	}
	if !user.IsAdmin && book.UserID.String() != userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not allowed"))
	}
	if err := bookService.DeleteBookWithData(req.Msg.BookId); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &api.DeleteBookResponse{Ok: true}
	return connect.NewResponse(resp), nil
}

func (s *Server) UpdateBookPage(
	ctx context.Context,
	req *connect.Request[api.UpdateBookPageRequest],
) (*connect.Response[api.UpdateBookPageResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
	}
	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	bookService := storage.NewBookService(storage.DB)
	if err := bookService.UpdateBookPage(req.Msg.BookId, userID, user.IsAdmin, req.Msg.Page); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &api.UpdateBookPageResponse{Ok: true}
	return connect.NewResponse(resp), nil
}
