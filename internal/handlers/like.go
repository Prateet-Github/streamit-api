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

type LikeHandler struct {
	likeRepo  *repositories.LikeRepository
	videoRepo *repositories.VideoRepository
}

func NewLikeHandler(
	likeRepo *repositories.LikeRepository,
	videoRepo *repositories.VideoRepository,
) *LikeHandler {
	return &LikeHandler{
		likeRepo:  likeRepo,
		videoRepo: videoRepo,
	}
}

func (h *LikeHandler) Like(c *gin.Context) {

	userID, videoID, ok := h.getUserAndVideo(c)
	if !ok {
		return
	}

	like, err := h.likeRepo.FindByUserAndVideo(
		c.Request.Context(),
		userID,
		videoID,
	)

	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if like != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Video already liked",
		})
		return
	}

	newLike := &models.Like{
		ID:        bson.NewObjectID(),
		UserID:    userID,
		VideoID:   videoID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.likeRepo.Create(
		c.Request.Context(),
		newLike,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to like video",
		})
		return
	}

	if err := h.videoRepo.IncrementLikes(
		c.Request.Context(),
		videoID,
		1,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update like count",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Video liked successfully",
	})
}

func (h *LikeHandler) Unlike(c *gin.Context) {

	userID, videoID, ok := h.getUserAndVideo(c)
	if !ok {
		return
	}

	like, err := h.likeRepo.FindByUserAndVideo(
		c.Request.Context(),
		userID,
		videoID,
	)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Like not found",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	if err := h.likeRepo.Delete(
		c.Request.Context(),
		like.ID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove like",
		})
		return
	}

	if err := h.videoRepo.IncrementLikes(
		c.Request.Context(),
		videoID,
		-1,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update like count",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Video unliked successfully",
	})
}

func (h *LikeHandler) getUserAndVideo(c *gin.Context) (
	bson.ObjectID,
	bson.ObjectID,
	bool,
) {

	videoID, err := bson.ObjectIDFromHex(c.Param("videoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid video ID",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	userIDValue, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	userIDHex, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	_, err = h.videoRepo.FindByID(
		c.Request.Context(),
		videoID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video not found",
		})
		return bson.ObjectID{}, bson.ObjectID{}, false
	}

	return userID, videoID, true
}
