package handlers

import (
	"net/http"
	"time"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubscriptionHandler struct {
	subRepo  *repositories.SubscriptionRepository
	userRepo *repositories.UserRepository
}

func NewSubscriptionHandler(
	subRepo *repositories.SubscriptionRepository,
	userRepo *repositories.UserRepository,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		subRepo:  subRepo,
		userRepo: userRepo,
	}
}

func (h *SubscriptionHandler) getSubscriberAndChannel(
	c *gin.Context,
) (bson.ObjectID, bson.ObjectID, bool) {

	channelID, err := bson.ObjectIDFromHex(c.Param("channelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid channel ID",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	userIDHex := c.GetString("userId")

	subscriberID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	_, err = h.userRepo.FindByID(
		c.Request.Context(),
		channelID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Channel not found",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	return subscriberID, channelID, true
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {

	subscriberID, channelID, ok := h.getSubscriberAndChannel(c)
	if !ok {
		return
	}

	if subscriberID == channelID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "You cannot subscribe to yourself",
		})
		return
	}

	sub, err := h.subRepo.FindBySubscriberAndChannel(
		c.Request.Context(),
		subscriberID,
		channelID,
	)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if sub != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Already subscribed",
		})
		return
	}

	newSub := &models.Subscription{
		ID:           bson.NewObjectID(),
		SubscriberID: subscriberID,
		ChannelID:    channelID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.subRepo.Create(
		c.Request.Context(),
		newSub,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to subscribe",
		})
		return
	}

	if err := h.userRepo.IncrementSubscribers(
		c.Request.Context(),
		channelID,
		1,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update subscriber count",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Subscribed successfully",
	})
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {

	subscriberID, channelID, ok := h.getSubscriberAndChannel(c)
	if !ok {
		return
	}

	subscription, err := h.subRepo.FindBySubscriberAndChannel(
		c.Request.Context(),
		subscriberID,
		channelID,
	)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Subscription not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if err := h.subRepo.Delete(
		c.Request.Context(),
		subscription.ID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to unsubscribe",
		})
		return
	}

	if err := h.userRepo.IncrementSubscribers(
		c.Request.Context(),
		channelID,
		-1,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update subscriber count",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unsubscribed successfully",
	})
}

func (h *SubscriptionHandler) GetSubscriptionStatus(c *gin.Context) {

	channelID, err := bson.ObjectIDFromHex(c.Param("channelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid channel ID",
		})
		return
	}

	channel, err := h.userRepo.FindByID(
		c.Request.Context(),
		channelID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Channel not found",
		})
		return
	}

	userIDHex := c.GetString("userId")

	subscribed := false

	if userIDHex != "" {

		subscriberID, err := bson.ObjectIDFromHex(userIDHex)
		if err == nil {

			subscription, err := h.subRepo.FindBySubscriberAndChannel(
				c.Request.Context(),
				subscriberID,
				channelID,
			)

			if err != nil && err != mongo.ErrNoDocuments {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				return
			}

			subscribed = subscription != nil
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"subscribed":       subscribed,
		"subscribersCount": channel.SubscribersCount,
	})
}

func (h *SubscriptionHandler) GetMySubscriptions(c *gin.Context) {

	userIDHex := c.GetString("userId")

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	subscriptions, err := h.subRepo.FindBySubscriber(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch subscriptions",
		})
		return
	}

	response := make([]models.UserResponse, 0, len(subscriptions))

	for _, sub := range subscriptions {

		channel, err := h.userRepo.FindByID(
			c.Request.Context(),
			sub.ChannelID,
		)
		if err != nil {
			continue
		}

		response = append(response, models.UserResponse{
			ID:               channel.ID.Hex(),
			Name:             channel.Name,
			Username:         channel.Username,
			Bio:              channel.Bio,
			SubscribersCount: channel.SubscribersCount,
		})
	}

	c.JSON(http.StatusOK, response)
}
