package delivery

import (
	"context"
	"errors"
	"log"
	"time"

	"connectrpc.com/connect"
	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func (s *Server) GetUser(
	ctx context.Context,
	req *connect.Request[api.GetUserRequest],
) (*connect.Response[api.GetUserResponse], error) {
	userID, err := shared.ValidateAccessToken(req)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("forbidden"))
	}

	if req.Msg.UserId != userID && !IsAdmin(userID) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("admin only"))
	}

	var user storage.User
	if err := storage.DB.First(&user, "id = ?", req.Msg.UserId).Error; err != nil {
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
	// Validate request
	if req.Msg.Username == "" || req.Msg.Password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username and password are required"))
	}
	if len(req.Msg.Username) < 5 && len(req.Msg.Password) < 5 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username must be at least 5 characters"))
	}
	// Create user in the database
	passwordHash, err := shared.HashPassword(req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if passwordHash == "" {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to hash password"))
	}

	user := storage.User{
		Username:     req.Msg.Username,
		PasswordHash: passwordHash,
		IsAdmin:      false,
	}

	if err := storage.DB.Create(&user).Error; err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
	if req.Msg.Username == "" || req.Msg.Password == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username and password are required"))
	}

	var user storage.User
	if err := storage.DB.First(&user, "username = ?", req.Msg.Username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("user not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !shared.CheckPassword(user.PasswordHash, req.Msg.Password) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid credentials"))
	}
	accessToken, err := GenerateAccessToken(user.ID.String())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate access token"))
	}
	refreshToken, err := GenerateRefreshToken(user.ID.String())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to generate refresh token"))
	}
	resp := &api.LoginResponse{
		Ok:           true,
		UserId:       user.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		IsAdmin:      user.IsAdmin,
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) RefreshToken(
	ctx context.Context,
	req *connect.Request[api.RefreshRequest],
) (*connect.Response[api.RefreshResponse], error) {
	log.Println("RefreshToken received:", req.Msg.RefreshToken)

	token, err := jwt.Parse(req.Msg.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return shared.JwtSecret, nil
	})
	if err != nil || !token.Valid {
		log.Println("Token is not valid:", err)
		resp := &api.RefreshResponse{
			AccessToken:  "",
			RefreshToken: "",
			Error:        "invalid refresh token",
		}
		return connect.NewResponse(resp), nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		log.Println("Claims error or wrong type:", claims)
		resp := &api.RefreshResponse{
			AccessToken:  "",
			RefreshToken: "",
			Error:        "invalid token type",
		}
		return connect.NewResponse(resp), nil
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Println("user_id not found in claims:", claims)
		resp := &api.RefreshResponse{
			AccessToken:  "",
			RefreshToken: "",
			Error:        "user_id not found",
		}
		return connect.NewResponse(resp), nil
	}
	accessToken, err := GenerateAccessToken(userID)
	if err != nil {
		log.Println("Failed to generate access token:", err)
		resp := &api.RefreshResponse{
			AccessToken:  "",
			RefreshToken: "",
			Error:        "failed to generate access token",
		}
		return connect.NewResponse(resp), nil
	}
	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		log.Println("Failed to generate refresh token:", err)
		resp := &api.RefreshResponse{
			AccessToken:  "",
			RefreshToken: "",
			Error:        "failed to generate refresh token",
		}
		return connect.NewResponse(resp), nil
	}
	resp := &api.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Error:        "",
	}
	log.Println("RefreshToken response:", resp)
	return connect.NewResponse(resp), nil
}

func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     jwt.NewNumericDate(time.Now().Add(shared.AccessTokenTTL)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(shared.JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(shared.JwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
