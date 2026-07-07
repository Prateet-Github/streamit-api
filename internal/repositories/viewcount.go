package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ViewRepository struct {
	collection *mongo.Collection
}

func NewViewRepository(db *database.Database) *ViewRepository {
	return &ViewRepository{
		collection: db.DB.Collection("videos"),
	}
}

func (r *ViewRepository) BulkIncrementViews(
	ctx context.Context,
	increments map[string]int64,
) error {

	if len(increments) == 0 {
		return nil
	}

	models := []mongo.WriteModel{}

	for videoID, count := range increments {
		objectID, err := bson.ObjectIDFromHex(videoID)
		if err != nil {
			return err
		}

		model := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": objectID}).
			SetUpdate(bson.M{
				"$inc": bson.M{
					"views": count,
				},
			})

		models = append(models, model)
	}

	_, err := r.collection.BulkWrite(ctx, models)
	return err

}
