package handlers

import (
	"strings"
	"time"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/models"
	s3util "github.com/Prateet-Github/streamit-api/internal/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	 "go.mongodb.org/mongo-driver/v2/bson"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type VideoHandler struct {
	s3Client *awss3.Client
	cfg      *config.Config
	videoRepo *repositories.VideoRepository
}

func NewVideoHandler(
	s3Client *awss3.Client,
	cfg *config.Config,
	videoRepo *repositories.VideoRepository,
) *VideoHandler {
	return &VideoHandler{
		s3Client:  s3Client,
		cfg:       cfg,
		videoRepo: videoRepo,
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

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

	video := &models.Video{
		ID:                   bson.NewObjectID(),
		Title:                req.Title,
		Description:          req.Description,
		S3Key:                req.S3Key,
		OwnerID:              ownerID,
		Status:               models.StatusPending,
		ProcessingProgress:   0,
		Views:                0,
		LikesCount:           0,
		Visibility:           models.VisibilityPublic,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	if err := h.videoRepo.Create(
		c.Request.Context(),
		video,
	); err != nil {
		c.JSON(500, gin.H{
			"error": "failed to create video",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "upload confirmed, video is pending processing",
		"video":   video,
	})
}
