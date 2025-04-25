package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ad struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Url         string             `bson:"url" json:"url"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
}

func (a *Ad) IsValid() error {
	if a.ID.Hex() == "" {
		return errors.New("Id is missing")
	}
	if a.Title == "" {
		return errors.New("Title is missing")
	}
	if a.Url == "" {
		return errors.New("Url is missing")
	}
	return nil
}
