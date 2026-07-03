package repositories

import (
	"context"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubscriptionRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionRepository(db *mongo.Database) *SubscriptionRepository {
	return &SubscriptionRepository{
		collection: db.Collection("subscriptions"),
	}
}

func (r *SubscriptionRepository) FindBySubscriberAndChannel(
	ctx context.Context,
	subscriberID bson.ObjectID,
	channelID bson.ObjectID,
) (*models.Subscription, error) {

	var subscription models.Subscription

	err := r.collection.FindOne(
		ctx,
		bson.M{
			"subscriberId": subscriberID,
			"channelId":    channelID,
		},
	).Decode(&subscription)

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) Create(
	ctx context.Context,
	subscription *models.Subscription,
) error {

	_, err := r.collection.InsertOne(
		ctx,
		subscription,
	)

	return err
}

func (r *SubscriptionRepository) Delete(
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

func (r *SubscriptionRepository) FindBySubscriber(
	ctx context.Context,
	subscriberID bson.ObjectID,
) ([]models.Subscription, error) {

	cursor, err := r.collection.Find(
		ctx,
		bson.M{
			"subscriberId": subscriberID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var subscriptions []models.Subscription

	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
