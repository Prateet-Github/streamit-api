package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
