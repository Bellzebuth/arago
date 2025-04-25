package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Bellzebuth/arago/tracker/internal/db"
	mygrpc "github.com/Bellzebuth/arago/tracker/internal/grpc"
	pb "github.com/Bellzebuth/arago/tracker/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	collection, err := db.InitMongo(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	port := os.Getenv("TRACKER_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := mygrpc.NewTrackerServer(collection)

	grpcServer := grpc.NewServer()
	pb.RegisterTrackerServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	fmt.Println(fmt.Sprintf("AdServer is running on port :%s...", port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
