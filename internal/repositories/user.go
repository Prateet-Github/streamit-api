package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		collection: db.DB.Collection("users"),
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {

	_, err := r.collection.InsertOne(ctx, user)

	return err
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	var user models.User

	err := r.collection.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUsername(

	ctx context.Context,
	username string,

) (*models.User, error) {

	var user models.User
	err := r.collection.FindOne(
		ctx,
		bson.M{"username": username},
	).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.User, error) {

	var user models.User

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateProfile(
	ctx context.Context,
	userID bson.ObjectID,
	name string,
	bio string,
) (*models.User, error) {

	update := bson.M{}

	if name != "" {
		update["name"] = name
	}

	if bio != "" {
		update["bio"] = bio
	}

	update["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var user models.User

	err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{
			"_id": userID,
		},
		bson.M{
			"$set": update,
		},
		opts,
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IncrementSubscribers(
	ctx context.Context,
	userID bson.ObjectID,
	delta int,
) error {

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id": userID,
		},
		bson.M{
			"$inc": bson.M{
				"subscribersCount": delta,
			},
		},
	)

	return err
}

func (r *VideoRepository) FindPublicByOwnerID(
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
			"status": models.StatusCompleted,
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