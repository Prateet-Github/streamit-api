package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateIndexes(db *Database) error {
	users := db.DB.Collection("users")

	_, err := users.Indexes().CreateMany(
		context.Background(),
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

	return err
}
