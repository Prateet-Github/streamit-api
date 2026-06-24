package handlers

import (
	"time"

	"errors"

	"github.com/Prateet-Github/streamit-api/internal/models"
	"github.com/Prateet-Github/streamit-api/internal/repositories"
	"github.com/Prateet-Github/streamit-api/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandler struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
}

func NewAuthHandler(
	userRepo *repositories.UserRepository,
	jwtSecret string,
) *AuthHandler {
	return &AuthHandler{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
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

	existingUser, err := h.userRepo.FindByEmail(
		c.Request.Context(),
		req.Email,
	)

	if err == nil && existingUser != nil {
		c.JSON(409, gin.H{
			"error": "email already exists",
		})
		return
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	existingUser, err = h.userRepo.FindByUsername(
		c.Request.Context(),
		req.Username,
	)

	if err == nil && existingUser != nil {
		c.JSON(409, gin.H{
			"error": "username already exists",
		})
		return
	}

	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(500, gin.H{
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userRepo.FindByEmail(
		c.Request.Context(),
		req.Email,
	)

	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	if !utils.CheckPassword(
		req.Password,
		user.Password,
	) {
		c.JSON(401, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	token, err := utils.GenerateToken(
		user.ID.Hex(),
		h.jwtSecret,
	)

	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}
