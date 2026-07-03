package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Subscription struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	SubscriberID bson.ObjectID `bson:"subscriberId" json:"subscriberId"`
	ChannelID    bson.ObjectID `bson:"channelId" json:"channelId"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
