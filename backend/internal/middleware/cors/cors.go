package cors

import (
    "github.com/rs/cors"
    "net/http"
)

func New() func(http.Handler) http.Handler {
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"*"},
        AllowCredentials: true,
    })
    return c.Handler
}