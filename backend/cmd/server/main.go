package main

import (
	"log"
	"net"

	"github.com/Paveluts42/bookreader/backend/api"
	"github.com/Paveluts42/bookreader/backend/internal/delivery"
	"github.com/Paveluts42/bookreader/backend/internal/storage"
	"google.golang.org/grpc"
)

func main() {
	storage.InitDB()
	storage.AutoMigrate()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	// TODO: register your service
	api.RegisterReaderServiceServer(grpcServer, &delivery.server{})

	log.Println("gRPC server listening on :50051")
  log.Println("gRPC server listening on :50051")
  grpcServer.Serve(lis)
}
