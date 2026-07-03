package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ChannelRepository struct {
	collection *mongo.Collection
}

func NewChannelRepository(db *mongo.Database) *ChannelRepository {
	return &ChannelRepository{
		collection: db.Collection("videos"),
	}
}

func (r *ChannelRepository) FindByOwnerID(
	ctx context.Context,
	ownerID bson.ObjectID,
) ([]models.Video, error) {

	opts := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		})

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"ownerId":    ownerID,
			"status":     models.StatusCompleted,
			"visibility": models.VisibilityPublic,
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var videos []models.Video

	if err := cursor.All(ctx, &videos); err != nil {
		return nil, err
	}

	return videos, nil
}
