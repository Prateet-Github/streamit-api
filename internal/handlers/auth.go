package handlers

import (
	"time"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
}

func NewAuthHandler(
	userRepo *repositories.UserRepository,
) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	user := &models.User{
		ID:        bson.NewObjectID(),
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		Bio:       "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.userRepo.Create(
		c.Request.Context(),
		user,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(201, user)
}
