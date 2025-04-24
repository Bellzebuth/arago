package grpc_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	adgrpc "github.com/Bellzebuth/arago/adserver/internal/grpc"
	pb "github.com/Bellzebuth/arago/adserver/proto/ad/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func initGRPCServer(t *testing.T, collection *mongo.Collection) pb.AdServiceClient {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	adServer := &adgrpc.AdServer{
		AdCollection: collection,
	}

	pb.RegisterAdServiceServer(s, adServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	// Create a gRPC client with grpc.NewClient
	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.Dial()
	}

	clientConn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	return pb.NewAdServiceClient(clientConn)
}

func setupMongo(t *testing.T) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Mongo connection error: %v", err)
	}

	collection := client.Database("ad_test").Collection("ads")
	// Clean up before each test
	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Could not drop test collection: %v", err)
	}

	return collection
}

func TestCreateAndGetAd(t *testing.T) {
	collection := setupMongo(t)
	client := initGRPCServer(t, collection)

	ctx := context.Background()

	// 1. Create Ad
	createResp, err := client.CreateAd(ctx, &pb.CreateAdRequest{
		Title:       "Test Ad",
		Description: "This is a test ad",
		Url:         "http://example.com",
	})
	if err != nil {
		t.Fatalf("CreateAd failed: %v", err)
	}

	if createResp.Ad.Title != "Test Ad" {
		t.Errorf("Expected title to be 'Test Ad', got %s", createResp.Ad.Title)
	}

	// 2. Get Ad
	getResp, err := client.GetAd(ctx, &pb.GetAdRequest{
		Id: createResp.Ad.Id,
	})
	if err != nil {
		t.Fatalf("GetAd failed: %v", err)
	}

	if getResp.Ad.Title != "Test Ad" {
		t.Errorf("Expected title to be 'Test Ad', got %s", getResp.Ad.Title)
	}
}
