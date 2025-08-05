package main

import (
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/delivery"
)

func main() {
	grpcServer := grpc.NewServer()
	api.RegisterReaderServiceServer(grpcServer, delivery.NewServer())

	wrapped := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return origin == "http://localhost:5173"
		}),
	)

	httpServer := &http.Server{
		Addr: ":50051",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Grpc-Web, X-User-Agent, Grpc-Timeout, Connect-Protocol-Version")
			w.Header().Set("Access-Control-Expose-Headers", "Grpc-Status, Grpc-Message, Grpc-Encoding, Grpc-Accept-Encoding")

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Прокидываем gRPC-Web запросы
				log.Println(">> grpc-web request:", r.URL.Path)
				wrapped.ServeHTTP(w, r)
	

			http.NotFound(w, r)
		}),
	}

	log.Println("Starting server on :50051")
	log.Fatal(httpServer.ListenAndServe())
}
