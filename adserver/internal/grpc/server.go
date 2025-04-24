package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Bellzebuth/arago/adserver/proto/ad/proto"
	trackerpb "github.com/Bellzebuth/arago/tracker/proto"
	"google.golang.org/grpc"

	"github.com/Bellzebuth/arago/adserver/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdServer struct {
	pb.UnimplementedAdServiceServer
	AdCollection  *mongo.Collection
	TrackerClient trackerpb.TrackerServiceClient
}

func (s *AdServer) initTrackerClient() error {
	conn, err := grpc.Dial("tracker:50052", grpc.WithInsecure()) // Assure-toi que le nom du service est correct
	if err != nil {
		return fmt.Errorf("failed to connect to tracker service: %v", err)
	}
	s.TrackerClient = trackerpb.NewTrackerServiceClient(conn)
	return nil
}

func (s *AdServer) CreateAd(ctx context.Context, req *pb.CreateAdRequest) (*pb.CreateAdResponse, error) {
	ad := models.Ad{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Url:         req.GetUrl(),
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour), // add expires after 24 hours
	}

	result, err := s.AdCollection.InsertOne(ctx, ad)
	if err != nil {
		return nil, err
	}

	ad.ID = result.InsertedID.(primitive.ObjectID)

	return &pb.CreateAdResponse{
		Ad: &pb.Ad{
			Id:          ad.ID.Hex(),
			Title:       ad.Title,
			Description: ad.Description,
			Url:         ad.Url,
		},
	}, nil
}

func (s *AdServer) GetAd(ctx context.Context, req *pb.GetAdRequest) (*pb.GetAdResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	var ad models.Ad
	err = s.AdCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ad)
	if err != nil {
		return nil, err
	}

	return &pb.GetAdResponse{
		Ad: &pb.Ad{
			Id:          ad.ID.Hex(),
			Title:       ad.Title,
			Description: ad.Description,
			Url:         ad.Url,
		},
	}, nil
}

func (s *AdServer) ServeAd(ctx context.Context, req *pb.ServeAdRequest) (*pb.ServeAdResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	var ad models.Ad
	err = s.AdCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ad)
	if err != nil {
		return nil, err
	}

	trackReq := &trackerpb.TrackClickRequest{
		AdId: ad.ID.Hex(),
		//Count: req.Count(),
	}

	_, err = s.TrackerClient.TrackClick(ctx, trackReq)
	if err != nil {
		log.Printf("Failed to track click: %v", err)
	}

	return &pb.ServeAdResponse{
		Url: ad.Url,
	}, nil
}
