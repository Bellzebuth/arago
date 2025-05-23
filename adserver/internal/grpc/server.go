package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Bellzebuth/arago/adserver/internal/cache"
	pb "github.com/Bellzebuth/arago/adserver/proto/proto"
	trackerpb "github.com/Bellzebuth/arago/tracker/proto/proto"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Bellzebuth/arago/adserver/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdServer struct {
	pb.UnimplementedAdServiceServer
	AdCollection  *mongo.Collection
	TrackerClient trackerpb.TrackerServiceClient
	RedisClient   *redis.Client
}

func (s *AdServer) Init(collection *mongo.Collection) error {
	trackerPort := os.Getenv("TRACKER_PORT")
	if trackerPort == "" {
		trackerPort = "50052"
	}

	conn, err := grpc.NewClient(fmt.Sprintf("tracker:%s", trackerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to tracker service: %v", err)
	}

	s.AdCollection = collection
	s.TrackerClient = trackerpb.NewTrackerServiceClient(conn)

	dragonflyPort := os.Getenv("DRAGONFLY_PORT")
	if dragonflyPort == "" {
		dragonflyPort = "6379"
	}
	s.RedisClient = cache.NewRedisClient(fmt.Sprintf("dragonfly:%s", dragonflyPort))
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

	err = ad.IsValid()
	if err != nil {
		return nil, err
	}

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
	cacheKey := fmt.Sprintf("ad:%s", req.GetId())

	// Try to get ad from cache
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedAd pb.Ad
		if err := json.Unmarshal([]byte(cached), &cachedAd); err == nil {
			return &pb.GetAdResponse{Ad: &cachedAd}, nil
		}
	}

	// If not in cache, query MongoDB
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	var ad models.Ad
	err = s.AdCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&ad)
	if err != nil {
		return nil, err
	}

	adResponse := &pb.Ad{
		Id:          ad.ID.Hex(),
		Title:       ad.Title,
		Description: ad.Description,
		Url:         ad.Url,
	}

	// Cache the result
	data, err := json.Marshal(adResponse)
	if err == nil {
		s.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute) // Set TTL as needed
	}

	return &pb.GetAdResponse{Ad: adResponse}, nil
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
	}

	_, err = s.TrackerClient.TrackClick(ctx, trackReq)
	if err != nil {
		log.Printf("Failed to track click: %v", err)
	}

	return &pb.ServeAdResponse{
		Url: ad.Url,
	}, nil
}
