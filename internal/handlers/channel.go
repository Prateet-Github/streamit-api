package handlers

import (
	"net/http"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	userRepo  *repositories.UserRepository
	videoRepo *repositories.VideoRepository
}

func NewChannelHandler(
	userRepo *repositories.UserRepository,
	videoRepo *repositories.VideoRepository,
) *ChannelHandler {
	return &ChannelHandler{
		userRepo:  userRepo,
		videoRepo: videoRepo,
	}
}

func (h *ChannelHandler) GetChannel(c *gin.Context) {

	username := c.Param("username")

	user, err := h.userRepo.FindByUsername(
		c.Request.Context(),
		username,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Channel not found",
		})
		return
	}

	videos, err := h.videoRepo.FindPublicByOwnerID(
		c.Request.Context(),
		user.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos",
		})
		return
	}

	videoResponse := make([]models.VideoListResponse, 0, len(videos))

	for _, video := range videos {

		videoResponse = append(videoResponse, models.VideoListResponse{
			ID:           video.ID.Hex(),
			Title:        video.Title,
			ThumbnailKey: video.ThumbnailKey,
			HLSURL:       video.HLSURL,
			CreatedAt:    video.CreatedAt,
			Views:        video.Views,
			Status:       video.Status,
			Owner: &models.OwnerResponse{
				ID:       user.ID.Hex(),
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user": models.UserResponse{
			ID:               user.ID.Hex(),
			Name:             user.Name,
			Username:         user.Username,
			Bio:              user.Bio,
			SubscribersCount: user.SubscribersCount,
		},
		"videos": videoResponse,
	})
}
