package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
				"likesCount": delta,
			},
		},
	)

	return err
}

// find all videos

func (r *VideoRepository) FindAll(
	ctx context.Context,
) ([]models.Video, error) {

	opts := options.Find().
		SetSort(
			bson.D{
				{Key: "createdAt", Value: -1},
			},
		)

	cursor, err := r.collection.Find(
		ctx,
		bson.M{},
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

// owner or video for the dashboard

func (r *VideoRepository) FindByOwnerID(
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
			"ownerId": ownerID,
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

func (r *VideoRepository) Search(
	ctx context.Context,
	query string,
) ([]models.Video, error) {

	filter := bson.M{
		"title": bson.M{
			"$regex":   query,
			"$options": "i",
		},
		"status":     models.StatusCompleted,
		"visibility": models.VisibilityPublic,
	}

	opts := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		}).
		SetLimit(20)

	cursor, err := r.collection.Find(
		ctx,
		filter,
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
