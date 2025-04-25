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

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "arago"
	}

	client, err := db.Connect(ctx, mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	adCollection, err := db.InitAdCollection(client.Database(dbName))
	if err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}

	adServer := adgrpc.AdServer{}
	err = adServer.Init(adCollection)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdServiceServer(grpcServer, &adServer)

	reflection.Register(grpcServer) // use to test with grpcurl

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Println(fmt.Sprintf("AdServer is running on port :%s...", port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
