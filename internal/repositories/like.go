package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LikeRepository struct {
	collection *mongo.Collection
}

func NewLikeRepository(db *mongo.Database) *LikeRepository {
	return &LikeRepository{
		collection: db.Collection("likes"),
	}
}

func (r *LikeRepository) FindByUserAndVideo(
	ctx context.Context,
	userID bson.ObjectID,
	videoID bson.ObjectID,
) (*models.Like, error) {

	filter := bson.M{
		"userId":  userID,
		"videoId": videoID,
	}

	var like models.Like

	err := r.collection.FindOne(ctx, filter).Decode(&like)
	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (r *LikeRepository) Create(
	ctx context.Context,
	like *models.Like,
) error {

	_, err := r.collection.InsertOne(ctx, like)
	return err
}

func (r *LikeRepository) Delete(
	ctx context.Context,
	id bson.ObjectID,
) error {

	_, err := r.collection.DeleteOne(
		ctx,
		bson.M{
			"_id": id,
		},
	)

	return err
}
