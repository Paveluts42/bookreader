package cors

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
)

// WithCORS adds CORS support to a Connect HTTP handler.
func WithCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
        AllowedOrigins: []string{
            "http://localhost:5173",
            "http://localhost:8090",
            "http://212.113.119.120:8090",
            "http://127.0.0.1:5173",
            "http://192.168.0.8:5173", 
        },
        AllowedMethods: connectcors.AllowedMethods(),
        AllowedHeaders: append(connectcors.AllowedHeaders(), "Authorization"),
        ExposedHeaders: append(connectcors.ExposedHeaders(), "Authorization"),
	})
	return middleware.Handler(h)
}