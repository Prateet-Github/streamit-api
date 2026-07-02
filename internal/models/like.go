package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Like struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"userId" json:"userId"`
	VideoID   bson.ObjectID `bson:"videoId" json:"videoId"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
