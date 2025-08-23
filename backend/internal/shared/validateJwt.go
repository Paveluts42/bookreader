package shared

import (
	"errors"

	"connectrpc.com/connect"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateAccessToken[T any](req *connect.Request[T]) (string, error) {
    tokenStr := req.Header().Get("Authorization")
    if tokenStr == "" {
        return "", errors.New("no authorization header")
    }
    if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
        tokenStr = tokenStr[7:]
    }
    if tokenStr == "" {
        return "", errors.New("no token provided")
    }

    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return JwtSecret, nil
    })
    if err != nil || !token.Valid {
        return "", errors.New("invalid token: " + err.Error())
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || claims["type"] != "access" {
        return "", errors.New("invalid token type")
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", errors.New("user_id not found in claims")
    }

    return userID, nil
}