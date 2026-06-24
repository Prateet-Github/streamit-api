package repositories

import (

	"context"
	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"

)


type VideoRepository struct {
	collection *mongo.Collection
}

func NewVideoRepository(db *database.Database) *VideoRepository {
	return &VideoRepository{
		collection: db.DB.Collection("videos"),
	}
}

func (r *VideoRepository) Create(
	ctx context.Context,
	video *models.Video,
) error {
	_, err := r.collection.InsertOne(ctx, video)
	return err
}