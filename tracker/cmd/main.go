package main

import (
	"log"
	"net"
	"os"

	"github.com/Bellzebuth/arago/tracker/internal/db"
	mygrpc "github.com/Bellzebuth/arago/tracker/internal/grpc"
	pb "github.com/Bellzebuth/arago/tracker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	collection, err := db.InitMongo(mongoURI, "arago", "clicks")
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := mygrpc.NewTrackerServer(collection)

	grpcServer := grpc.NewServer()
	pb.RegisterTrackerServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	log.Println("Tracker service is running on port 50052...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
