package models

import "time"

type Click struct {
	AdID      string    `bson:"ad_id"`
	Timestamp time.Time `bson:"timestamp"`
}
