package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (r *VideoRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Video, error) {

	var video models.Video

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&video)

	if err != nil {
		return nil, err
	}

	return &video, nil
}

func (r *VideoRepository) Update(
	ctx context.Context,
	video *models.Video,
) error {

	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{
			"_id": video.ID,
		},
		video,
	)

	return err
}

func (r *VideoRepository) IncrementLikes(
	ctx context.Context,
	videoID bson.ObjectID,
	delta int,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": videoID,
		},
		bson.M{
			"$inc": bson.M{
				"likes": delta,
			},
		},
	)

	return err
}
