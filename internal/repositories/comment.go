package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CommentRepository struct {
	collection *mongo.Collection
}

func NewCommentRepository(db *mongo.Database) *CommentRepository {
	return &CommentRepository{
		collection: db.Collection("comments"),
	}
}

func (r *CommentRepository) Create(
	ctx context.Context,
	comment *models.Comment,
) error {

	_, err := r.collection.InsertOne(
		ctx,
		comment,
	)

	return err
}

func (r *CommentRepository) FindByVideoID(
	ctx context.Context,
	videoID bson.ObjectID,
) ([]models.Comment, error) {

	opts := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: -1},
		})

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"videoId": videoID,
			"parentCommentId": bson.M{
				"$exists": false,
			},
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.Comment

	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.Comment, error) {

	var comment models.Comment

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"_id": id,
		},
	).Decode(&comment)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *CommentRepository) Delete(
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

func (r *CommentRepository) FindReplies(
	ctx context.Context,
	parentCommentID bson.ObjectID,
) ([]models.Comment, error) {

	opts := options.Find().
		SetSort(bson.D{
			{Key: "createdAt", Value: 1},
		})

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"parentCommentId": parentCommentID,
		},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var replies []models.Comment

	if err := cursor.All(ctx, &replies); err != nil {
		return nil, err
	}

	return replies, nil
}

func (r *CommentRepository) IncrementRepliesCount(
	ctx context.Context,
	commentID bson.ObjectID,
	delta int,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": commentID,
		},
		bson.M{
			"$inc": bson.M{
				"repliesCount": delta,
			},
		},
	)

	return err
}

func (r *VideoRepository) IncrementCommentsCount(
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
				"commentsCount": delta,
			},
		},
	)

	return err
}
