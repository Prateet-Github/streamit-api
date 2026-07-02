package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateIndexes(db *Database) error {
	ctx := context.Background()

	users := db.DB.Collection("users")
	likes := db.DB.Collection("likes")

	_, err := users.Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys: bson.M{"email": 1},
				Options: options.Index().
					SetUnique(true),
			},
			{
				Keys: bson.M{"username": 1},
				Options: options.Index().
					SetUnique(true),
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = likes.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "videoId", Value: 1},
			},
			Options: options.Index().
				SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
