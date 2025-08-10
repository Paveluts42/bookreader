package main

import (
	"log"
	"net/http"

	apiconnect "github.com/Paveluts42/bookreader/backend/api/apiconnect"
	"github.com/Paveluts42/bookreader/backend/internal/delivery"
	corsmw "github.com/Paveluts42/bookreader/backend/internal/middleware/cors"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
)

func main() {
	storage.InitDB()
	storage.AutoMigrate(storage.DB)
	mux := http.NewServeMux()

	server := delivery.NewServer()
	path, handler := apiconnect.NewReaderServiceHandler(server)

	mux.Handle(path, handler)
	mux.Handle("/uploads/", corsmw.New()(http.StripPrefix("/uploads/", http.FileServer(http.Dir("/uploads")))))
	// Add CORS middleware
	handlerWithCORS := corsmw.New()(mux)
	log.Println("🚀 Starting CONNECT server on :50051")
	if err := http.ListenAndServe(":50051", handlerWithCORS); err != nil {
		log.Fatal(err)
	}
}
