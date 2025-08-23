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
	bookPath, bookHandler := apiconnect.NewBookServiceHandler(server)
	notePath, noteHandler := apiconnect.NewNoteServiceHandler(server)
	userPath, userHandler := apiconnect.NewUserServiceHandler(server)

	mux.Handle(bookPath, bookHandler)
	mux.Handle(notePath, noteHandler)
	mux.Handle(userPath, userHandler)
	mux.Handle("/uploads/", corsmw.WithCORS(http.StripPrefix("/uploads/", http.FileServer(http.Dir("/uploads")))))
	// Add CORS middleware
	handlerWithCORS := corsmw.WithCORS(mux)
	log.Println("🚀 Starting CONNECT server on :50051")
	if err := http.ListenAndServe(":50051", handlerWithCORS); err != nil {
		log.Fatal(err)
	}
}
