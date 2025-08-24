package delivery

import (
    "context"
    "errors"

    "connectrpc.com/connect"
    "github.com/Paveluts42/bookreader/backend/api"
    "github.com/Paveluts42/bookreader/backend/internal/shared"
    "github.com/Paveluts42/bookreader/backend/internal/storage"
)

func (s *Server) GetUser(
    ctx context.Context,
    req *connect.Request[api.GetUserRequest],
) (*connect.Response[api.GetUserResponse], error) {
    userID, err := shared.ValidateAccessToken(req)
    if err != nil {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
    }
    authService := storage.NewAuthService(storage.DB)
    if req.Msg.UserId != userID && !authService.IsAdmin(userID) {
        return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
    }
    user, err := authService.GetUserByID(req.Msg.UserId)
    if err != nil {
        return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
    }
    resp := &api.GetUserResponse{
        UserId:   user.ID.String(),
        Username: user.Username,
        IsAdmin:  user.IsAdmin,
    }
    return connect.NewResponse(resp), nil
}

func (s *Server) Register(
    ctx context.Context,
    req *connect.Request[api.RegisterRequest],
) (*connect.Response[api.RegisterResponse], error) {
    authService := storage.NewAuthService(storage.DB)
    user, err := authService.Register(req.Msg.Username, req.Msg.Password)
    if err != nil {
        return nil, connect.NewError(connect.CodeInvalidArgument, err)
    }
    resp := &api.RegisterResponse{
        Ok:     true,
        UserId: user.ID.String(),
    }
    return connect.NewResponse(resp), nil
}

func (s *Server) Login(
    ctx context.Context,
    req *connect.Request[api.LoginRequest],
) (*connect.Response[api.LoginResponse], error) {
    authService := storage.NewAuthService(storage.DB)
    loginResp, err := authService.Login(req.Msg.Username, req.Msg.Password)
    if err != nil {
        return nil, connect.NewError(connect.CodeUnauthenticated, err)
    }
    return connect.NewResponse(loginResp), nil
}

func (s *Server) RefreshToken(
    ctx context.Context,
    req *connect.Request[api.RefreshRequest],
) (*connect.Response[api.RefreshResponse], error) {
    authService := storage.NewAuthService(storage.DB)
    resp, err := authService.RefreshToken(req.Msg.RefreshToken)
    if err != nil {
        return connect.NewResponse(resp), nil
    }
    return connect.NewResponse(resp), nil
}