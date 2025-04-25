package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Bellzebuth/arago/tracker/models"
	pb "github.com/Bellzebuth/arago/tracker/proto"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrackerServer struct {
	pb.UnimplementedTrackerServiceServer
	ClickCollection *mongo.Collection
}

func NewTrackerServer(collection *mongo.Collection) *TrackerServer {
	return &TrackerServer{
		ClickCollection: collection,
	}
}

func (s *TrackerServer) TrackClick(ctx context.Context, req *pb.TrackClickRequest) (*pb.TrackClickResponse, error) {
	click := models.Click{
		AdID:      req.GetAdId(),
		Timestamp: time.Now(),
	}

	fmt.Println("LAAAAAAAAAAAAAAAAAAAA")
	_, err := s.ClickCollection.InsertOne(ctx, click)
	if err != nil {
		log.Printf("Failed to insert click: %v", err)
		return &pb.TrackClickResponse{Success: false}, err
	}
	fmt.Println("LOOOOOOOOOOOOOOOOO")

	return &pb.TrackClickResponse{Success: true}, nil
}
