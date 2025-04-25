package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Database("admin").RunCommand(ctx, map[string]any{"ping": 1}).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func InitAdCollection(db *mongo.Database) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collectionName := os.Getenv("AD_COLLECTION")
	if collectionName == "" {
		collectionName = "ads"
	}

	collection := db.Collection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

	return collection, nil
}
