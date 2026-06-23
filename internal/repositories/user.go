package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"
	"github.com/Prateet-Github/streamit-api/internal/models"
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
