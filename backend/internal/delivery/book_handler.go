package delivery

import (
	"context"
	"errors"
	"log"
	"os"

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


	if err != nil {
		log.Println("Error:", err.Error())
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}
	if userID == "" {
		log.Println("Error: userID is empty")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("userID is empty"))
	}

    var books []storage.Book
    var user storage.User
    if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
        log.Printf("DB error: %v", err)
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    if user.IsAdmin {
        // Админ получает все книги
        if err := storage.DB.Find(&books).Error; err != nil {
            log.Printf("DB error: %v", err)
            return nil, connect.NewError(connect.CodeInternal, err)
        }
    } else {
        // Обычный пользователь — только свои
        if err := storage.DB.Where("user_id = ?", userID).Find(&books).Error; err != nil {
            log.Printf("DB error: %v", err)
            return nil, connect.NewError(connect.CodeInternal, err)
        }
    }

	resp := &api.GetBooksResponse{
		Books: make([]*api.Book, 0, len(books)),
	}
	for _, b := range books {
		resp.Books = append(resp.Books, &api.Book{
			Id:       b.ID.String(),
			Title:    b.Title,
			Author:   b.Author,
			FilePath: b.FilePath,
			CoverUrl: b.CoverPath,
			UserId: b.UserID.String(),
			Page:     int32(b.Page),
			PageAll:  int32(b.PageAll),
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

    var book storage.Book
    if user.IsAdmin {
        err = storage.DB.First(&book, "id = ?", req.Msg.BookId).Error
    } else {
        err = storage.DB.Where("id = ? AND user_id = ?", req.Msg.BookId, userID).First(&book).Error
    }

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, connect.NewError(connect.CodeNotFound, err)
        }
        log.Printf("DB error: %v", err)
        return nil, connect.NewError(connect.CodeInternal, err)
    }


    bookResp := &api.Book{
        Id:       book.ID.String(),
        Title:    book.Title,
        Author:   book.Author,
        FilePath: book.FilePath,
        CoverUrl: book.CoverPath,
        Page:     int32(book.Page),
        PageAll:  int32(book.PageAll),
    }
    resp := &api.GetBookResponse{
        Book: bookResp,
    }
    return connect.NewResponse(resp), nil
}

func (s *Server) DeleteBook(
	ctx context.Context,
	req *connect.Request[api.DeleteBookRequest],
) (*connect.Response[api.DeleteBookResponse], error) {
	var book storage.Book
	if err := storage.DB.First(&book, "id = ?", req.Msg.BookId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		log.Printf("DB error: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := storage.DB.Where("book_id = ?", book.ID).Delete(&storage.Note{}).Error; err != nil {
		log.Printf("Failed to delete notes: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := storage.DB.Delete(&book).Error; err != nil {
		log.Printf("Failed to delete book: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := os.Remove(book.FilePath); err != nil {
		log.Printf("Failed to delete file: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := os.Remove(book.CoverPath); err != nil {
		log.Printf("Failed to delete cover file: %v", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.DeleteBookResponse{Ok: true}
	return connect.NewResponse(resp), nil
}

func (s *Server) UpdateBookPage(
	ctx context.Context,
	req *connect.Request[api.UpdateBookPageRequest],
) (*connect.Response[api.UpdateBookPageResponse], error) {
	bookID := req.Msg.BookId
	page := req.Msg.Page

	var book storage.Book
	if err := storage.DB.First(&book, "id = ?", bookID).Error; err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	book.Page = page
	if err := storage.DB.Save(&book).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.UpdateBookPageResponse{Ok: true}
	return connect.NewResponse(resp), nil
}
