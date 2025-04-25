package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(ctx context.Context, uri string) (*mongo.Collection, error) {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "arago"
	}

	collectionName := os.Getenv("TARGET_COLLECTION")
	if collectionName == "" {
		collectionName = "tracker"
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return collection, nil
}
