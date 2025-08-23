package cors

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
)

// WithCORS adds CORS support to a Connect HTTP handler.
func WithCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: connectcors.AllowedMethods(),
        AllowedHeaders: append(connectcors.AllowedHeaders(), "Authorization"),
        ExposedHeaders: append(connectcors.ExposedHeaders(), "Authorization"),
	})
	return middleware.Handler(h)
}