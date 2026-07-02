package handlers

import (
	"log"
	"strings"
	"time"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	s3util "github.com/Prateet-Github/streamit-api/internal/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"

	"net/http"

	"github.com/Prateet-Github/streamit-api/internal/queue"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hibiken/asynq"
)

type VideoHandler struct {
	s3Client    *awss3.Client
	cfg         *config.Config
	videoRepo   *repositories.VideoRepository
	asynqClient *asynq.Client
	userRepo  *repositories.UserRepository
}

func NewVideoHandler(
	s3Client *awss3.Client,
	cfg *config.Config,
	videoRepo *repositories.VideoRepository,
	userRepo *repositories.UserRepository,
	asynqClient *asynq.Client,
) *VideoHandler {
	return &VideoHandler{
		s3Client:    s3Client,
		cfg:         cfg,
		videoRepo:   videoRepo,
		asynqClient: asynqClient,
		userRepo:  userRepo,
	}
}

func (h *VideoHandler) GetUploadURL(c *gin.Context) {
	var req models.UploadURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !strings.HasPrefix(req.FileType, "video/") {
		c.JSON(400, gin.H{
			"error": "only video files allowed",
		})
		return
	}

	key := "uploads/" + uuid.NewString()

	uploadURL, err := s3util.GenerateUploadURL(
		h.s3Client,
		h.cfg.S3RawBucket,
		key,
		req.FileType,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate upload url",
		})
		return
	}

	c.JSON(200, gin.H{
		"uploadUrl": uploadURL,
		"key":       key,
	})
}

func (h *VideoHandler) ConfirmUpload(c *gin.Context) {
	var req models.ConfirmUploadRequest

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if title, description and s3Key are provided
	if req.Title == "" || req.Description == "" || req.S3Key == "" {
		c.JSON(400, gin.H{
			"error": "title, description and s3Key are required",
		})
		return
	}

	userID := c.GetString("userId")

	ownerID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// create a new video
	video := &models.Video{
		ID:                 bson.NewObjectID(),
		Title:              req.Title,
		Description:        req.Description,
		S3Key:              req.S3Key,
		OwnerID:            ownerID,
		Status:             models.StatusPending,
		ProcessingProgress: 0,
		Views:              0,
		LikesCount:         0,
		Visibility:         models.VisibilityPublic,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// save to db
	if err := h.videoRepo.Create(
		c.Request.Context(),
		video,
	); err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create video",
		})
		return
	}

	// enqueue a task to process the video
	task, err := queue.NewProcessVideoTask(
		queue.VideoTask{
			VideoID: video.ID.Hex(),
			S3Key:   video.S3Key,
		},
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create processing task",
		})
		return
	}

	info, err := h.asynqClient.Enqueue(task)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to enqueue processing task",
		})
		return
	}

	log.Printf(
		"Enqueued task %s for video %s",
		info.ID,
		video.ID.Hex(),
	)

	// return the video object to the client
	c.JSON(201, gin.H{
		"message": "upload confirmed, video is pending processing",
		"video":   video,
	})
}

func (h *VideoHandler) CompleteVideo(c *gin.Context) {

	videoID := c.Param("id")

	var req models.CompleteVideoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := bson.ObjectIDFromHex(videoID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid video id",
		})
		return
	}

	video, err := h.videoRepo.FindByID(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "video not found",
		})
		return
	}

	video.Status = models.StatusCompleted
	video.HLSURL = req.HLSURL
	video.ThumbnailKey = req.ThumbnailKey
	video.ProcessingProgress = 100
	video.UpdatedAt = time.Now()

	if err := h.videoRepo.Update(
		c.Request.Context(),
		video,
	); err != nil {
		c.JSON(500, gin.H{
			"error": "failed to update video",
		})
		return
	}

	c.Status(204)
}

func (h *VideoHandler) GetAllVideos(c *gin.Context) {

videos, err := h.videoRepo.FindAll(
	c.Request.Context(),
)
if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Failed to fetch videos",
	})
	return
}

response := make([]models.VideoListResponse, 0, len(videos))

for _, video := range videos {

	user, err := h.userRepo.FindByID(
		c.Request.Context(),
		video.OwnerID,
	)

	if err != nil {
		continue
	}

	response = append(response, models.VideoListResponse{
		ID:           video.ID.Hex(),
		Title:        video.Title,
		ThumbnailKey: video.ThumbnailKey,
		HLSURL:       video.HLSURL,
		CreatedAt:    video.CreatedAt,
		Views:        video.Views,
		Owner: &models.OwnerResponse{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	})
}

c.JSON(http.StatusOK, response)
}
// single whole video

func (h *VideoHandler) GetVideoByID(c *gin.Context) {

	videoID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid video ID",
		})
		return
	}

	video, err := h.videoRepo.FindByID(
		c.Request.Context(),
		videoID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video not found",
		})
		return
	}

	user, err := h.userRepo.FindByID(
		c.Request.Context(),
		video.OwnerID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Owner not found",
		})
		return
	}

	response := models.VideoResponse{
		ID:           video.ID.Hex(),
		Title:        video.Title,
		Description:  video.Description,
		HLSURL:       video.HLSURL,
		ThumbnailKey: video.ThumbnailKey,
		Status:       video.Status,
		Owner: &models.OwnerResponse{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *VideoHandler) GetMyVideos(c *gin.Context) {

	userIDHex := c.GetString("userId")

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	videos, err := h.videoRepo.FindByOwnerID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch videos",
		})
		return
	}

	response := make([]models.VideoListResponse, 0, len(videos))

	user, err := h.userRepo.FindByID(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "User not found",
		})
		return
	}

	for _, video := range videos {

		response = append(response, models.VideoListResponse{
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

	c.JSON(http.StatusOK, response)
}