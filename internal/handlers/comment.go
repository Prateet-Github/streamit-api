package handlers

import (
	"net/http"
	"time"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CommentHandler struct {
	commentRepo *repositories.CommentRepository
	videoRepo   *repositories.VideoRepository
	userRepo    *repositories.UserRepository
}

func NewCommentHandler(
	commentRepo *repositories.CommentRepository,
	videoRepo *repositories.VideoRepository,
	userRepo *repositories.UserRepository,
) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
		videoRepo:   videoRepo,
		userRepo:    userRepo,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {

	videoID, err := bson.ObjectIDFromHex(c.Param("videoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid video ID",
		})
		return
	}

	_, err = h.videoRepo.FindByID(
		c.Request.Context(),
		videoID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Video not found",
		})
		return
	}

	userIDHex := c.GetString("userId")

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	var req models.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	comment := &models.Comment{
		ID:           bson.NewObjectID(),
		UserID:       userID,
		VideoID:      videoID,
		Content:      req.Content,
		RepliesCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.commentRepo.Create(
		c.Request.Context(),
		comment,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create comment",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Comment created",
	})
}

func (h *CommentHandler) GetComments(c *gin.Context) {

	videoID, err := bson.ObjectIDFromHex(c.Param("videoId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid video ID",
		})
		return
	}

	comments, err := h.commentRepo.FindByVideoID(
		c.Request.Context(),
		videoID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch comments",
		})
		return
	}

	response := make([]models.CommentResponse, 0, len(comments))

	for _, comment := range comments {

		user, err := h.userRepo.FindByID(
			c.Request.Context(),
			comment.UserID,
		)
		if err != nil {
			continue
		}

		response = append(response, models.CommentResponse{
			ID:           comment.ID.Hex(),
			Content:      comment.Content,
			RepliesCount: comment.RepliesCount,
			CreatedAt:    comment.CreatedAt,
			User: &models.OwnerResponse{
				ID:       user.ID.Hex(),
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {

	commentID, err := bson.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	comment, err := h.commentRepo.FindByID(
		c.Request.Context(),
		commentID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Comment not found",
		})
		return
	}

	userIDHex := c.GetString("userId")

	userID, _ := bson.ObjectIDFromHex(userIDHex)

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	if err := h.commentRepo.Delete(
		c.Request.Context(),
		commentID,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted",
	})
}

func (h *CommentHandler) CreateReply(c *gin.Context) {

	parentID, err := bson.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	parent, err := h.commentRepo.FindByID(
		c.Request.Context(),
		parentID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Parent comment not found",
		})
		return
	}

	userIDHex := c.GetString("userId")

	userID, err := bson.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user",
		})
		return
	}

	var req models.CreateCommentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	reply := &models.Comment{
		ID:              bson.NewObjectID(),
		UserID:          userID,
		VideoID:         parent.VideoID,
		Content:         req.Content,
		ParentCommentID: &parentID,
		RepliesCount:    0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.commentRepo.Create(
		c.Request.Context(),
		reply,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create reply",
		})
		return
	}

	if err := h.commentRepo.IncrementRepliesCount(
		c.Request.Context(),
		parentID,
		1,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update reply count",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reply created",
	})
}

func (h *CommentHandler) GetReplies(c *gin.Context) {

	commentID, err := bson.ObjectIDFromHex(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		return
	}

	replies, err := h.commentRepo.FindReplies(
		c.Request.Context(),
		commentID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch replies",
		})
		return
	}

	response := make([]models.CommentResponse, 0, len(replies))

	for _, reply := range replies {

		user, err := h.userRepo.FindByID(
			c.Request.Context(),
			reply.UserID,
		)
		if err != nil {
			continue
		}

		response = append(response, models.CommentResponse{
			ID:           reply.ID.Hex(),
			Content:      reply.Content,
			RepliesCount: reply.RepliesCount,
			CreatedAt:    reply.CreatedAt,
			User: &models.OwnerResponse{
				ID:       user.ID.Hex(),
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}
