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
	// bind request body to struct
	var req models.RegisterRequest

	// validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check if email already exists
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

	// check if username already exists
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

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// create user
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

	// return user
	c.JSON(201, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	// bind request body to struct
	var req models.LoginRequest

	// validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// find user by email
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

	// check if password is correct
	if !utils.CheckPassword(
		req.Password,
		user.Password,
	) {
		c.JSON(401, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// generate token
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

	// return token and user
	c.JSON(200, gin.H{
		"message": "login successful",
		"token":   token,
		"user":    user,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	// get user id from context
	userID, exists := c.Get("userId")

	if !exists {
		c.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// convert user id to bson.ObjectID
	objectID, err := bson.ObjectIDFromHex(
		userID.(string),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid user id",
		})
		return
	}

	// find user by id
	user, err := h.userRepo.FindByID(
		c.Request.Context(),
		objectID,
	)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	// return user
	c.JSON(200, gin.H{
		"user": user,
	})
}
