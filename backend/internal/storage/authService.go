package storage

import (
	"errors"
	"log"
	"time"

	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthService struct {
    db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
    return &AuthService{db: db}
}

func (s *AuthService) IsAdmin(userID string) bool {
    var user User
    if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
        return false
    }
    return user.IsAdmin
}

func (s *AuthService) GetUserByID(userID string) (*User, error) {
    var user User
    if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *AuthService) Register(username, password string) (*User, error) {
    if username == "" || password == "" {
        return nil, errors.New("username and password are required")
    }
    if len(username) < 5 || len(password) < 5 {
        return nil, errors.New("username and password must be at least 5 characters")
    }
    passwordHash, err := shared.HashPassword(password)
    if err != nil || passwordHash == "" {
        return nil, errors.New("failed to hash password")
    }
    user := User{
        Username:     username,
        PasswordHash: passwordHash,
        IsAdmin:      false,
    }
    if err := s.db.Create(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (s *AuthService) Login(username, password string) (*api.LoginResponse, error) {
    var user User
    if err := s.db.First(&user, "username = ?", username).Error; err != nil {
        return nil, errors.New("user not found")
    }
    if !shared.CheckPassword(user.PasswordHash, password) {
        return nil, errors.New("invalid credentials")
    }
    accessToken, err := GenerateAccessToken(user.ID.String())
    if err != nil {
        return nil, errors.New("failed to generate access token")
    }
    refreshToken, err := GenerateRefreshToken(user.ID.String())
    if err != nil {
        return nil, errors.New("failed to generate refresh token")
    }
    return &api.LoginResponse{
        Ok:           true,
        UserId:       user.ID.String(),
        Username:     user.Username,
        RefreshToken: refreshToken,
        AccessToken:  accessToken,
        IsAdmin:      user.IsAdmin,
    }, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*api.RefreshResponse, error) {
    log.Println("RefreshToken received:", refreshToken)
    token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
        return shared.JwtSecret, nil
    })
    if err != nil || !token.Valid {
        return &api.RefreshResponse{Error: "invalid refresh token"}, nil
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["type"] != "refresh" {
        return &api.RefreshResponse{Error: "invalid token type"}, nil
    }
    userID, ok := claims["user_id"].(string)
    if !ok {
        return &api.RefreshResponse{Error: "user_id not found"}, nil
    }
    accessToken, err := GenerateAccessToken(userID)
    if err != nil {
        return &api.RefreshResponse{Error: "failed to generate access token"}, nil
    }
    newRefreshToken, err := GenerateRefreshToken(userID)
    if err != nil {
        return &api.RefreshResponse{Error: "failed to generate refresh token"}, nil
    }
    return &api.RefreshResponse{
        AccessToken:  accessToken,
        RefreshToken: newRefreshToken,
        Error:        "",
    }, nil
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
        "exp":     jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(shared.JwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}