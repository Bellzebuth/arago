package models

type Click struct {
	AdID  string `bson:"ad_id"`
	Count int    `bson:"count"`
}
