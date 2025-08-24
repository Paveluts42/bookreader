package delivery

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
)



func (s *Server) GetUsers(
	ctx context.Context,
	req *connect.Request[api.GetUsersRequest],
) (*connect.Response[api.GetUsersResponse], error) {
    userID, err := shared.ValidateAccessToken(req)
    userService := storage.NewUserService(storage.DB)
    if err != nil || !userService.IsAdmin(userID) {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
    }
    users, err := userService.GetUsers()
    if err != nil {
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
	userService := storage.NewUserService(storage.DB)
	if err != nil || !userService.IsAdmin(userID) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
	}
	if req.Msg.UserId == userID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("cannot delete yourself"))
	}
	if err := userService.DeleteUserWithData(req.Msg.UserId); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	resp := &api.DeleteUserResponse{Ok: true}
	return connect.NewResponse(resp), nil
}
