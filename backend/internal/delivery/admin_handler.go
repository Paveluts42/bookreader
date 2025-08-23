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
)

func IsAdmin(userID string) bool {
	var user storage.User
	if err := storage.DB.First(&user, "id = ?", userID).Error; err != nil {
		return false
	}
	return user.IsAdmin
}

func (s *Server) GetUsers(
	ctx context.Context,
	req *connect.Request[api.GetUsersRequest],
) (*connect.Response[api.GetUsersResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil || !IsAdmin(userID) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
	}

	var users []storage.User
	if err := storage.DB.Find(&users).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &api.GetUsersResponse{
		Users: make([]*api.User, 0, len(users)),
	}
	for _, u := range users {
		resp.Users = append(resp.Users, &api.User{
			Id:       u.ID.String(),
			Username: u.Username,
			IsAdmin:  u.IsAdmin,
		})
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) DeleteUser(
	ctx context.Context,
	req *connect.Request[api.DeleteUserRequest],
) (*connect.Response[api.DeleteUserResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil || !IsAdmin(userID) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
	}
	if req.Msg.UserId == userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("cannot delete yourself"))
	}
	var books []storage.Book
	if err := storage.DB.Where("user_id = ?", req.Msg.UserId).Find(&books).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	for _, book := range books {
		if err := storage.DB.Where("book_id = ?", book.ID).Delete(&storage.Note{}).Error; err != nil {
			log.Printf("Failed to delete notes for book %s: %v", book.ID, err)
		}
		if err := storage.DB.Delete(&book).Error; err != nil {
			log.Printf("Failed to delete book %s: %v", book.ID, err)
		}
		if book.FilePath != "" {
			if err := os.Remove(book.FilePath); err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to delete PDF file %s: %v", book.FilePath, err)
			}
		}
		if book.CoverPath != "" {
			if err := os.Remove(book.CoverPath); err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to delete PNG file %s: %v", book.CoverPath, err)
			}
		}
	}

	if err := storage.DB.Delete(&storage.User{}, "id = ?", req.Msg.UserId).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &api.DeleteUserResponse{Ok: true}
	return connect.NewResponse(resp), nil
}
