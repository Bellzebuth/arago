package grpc

import (
	"context"
	"time"

	pb "github.com/Bellzebuth/arago/adserver/proto/ad/proto"

	"github.com/Bellzebuth/arago/adserver/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdServer struct {
	pb.UnimplementedAdServiceServer
	AdCollection *mongo.Collection
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

	// TODO: send gRPC request to tracker

	return &pb.ServeAdResponse{
		Url: ad.Url,
	}, nil
}
