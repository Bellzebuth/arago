package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/Bellzebuth/arago/adserver/internal/db"
	adgrpc "github.com/Bellzebuth/arago/adserver/internal/grpc"
	pb "github.com/Bellzebuth/arago/adserver/proto/ad/proto"

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

	client, err := db.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	adCollection, err := db.InitAdCollection(client.Database("adserver"))
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}

	s := adgrpc.AdServer{
		AdCollection: adCollection,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdServiceServer(grpcServer, &s)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println("AdServer is running on port :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
