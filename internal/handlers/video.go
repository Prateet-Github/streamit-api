package handlers

import (
	"strings"

	"github.com/Prateet-Github/streamit-api/internal/config"
	"github.com/Prateet-Github/streamit-api/internal/models"
	s3util "github.com/Prateet-Github/streamit-api/internal/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type VideoHandler struct {
	s3Client *awss3.Client
	cfg      *config.Config
}

func NewVideoHandler(
	s3Client *awss3.Client,
	cfg *config.Config,
) *VideoHandler {
	return &VideoHandler{
		s3Client: s3Client,
		cfg:      cfg,
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
