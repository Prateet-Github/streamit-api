package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type VideoStatus string

const (
	StatusPending    VideoStatus = "PENDING"
	StatusProcessing VideoStatus = "PROCESSING"
	StatusCompleted  VideoStatus = "COMPLETED"
	StatusFailed     VideoStatus = "FAILED"
)

type Visibility string

const (
	VisibilityPublic  Visibility = "PUBLIC"
	VisibilityPrivate Visibility = "PRIVATE"
)

type Video struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`

	S3Key string `bson:"s3Key" json:"s3Key"`

	HLSURL string `bson:"hlsUrl,omitempty" json:"hlsUrl,omitempty"`

	ThumbnailKey string `bson:"thumbnailKey,omitempty" json:"thumbnailKey,omitempty"`

	OwnerID bson.ObjectID `bson:"ownerId" json:"ownerId"`

	Status VideoStatus `bson:"status" json:"status"`

	ProcessingProgress int `bson:"processingProgress" json:"processingProgress"`

	ErrorMessage string `bson:"errorMessage,omitempty" json:"errorMessage,omitempty"`

	Duration int `bson:"duration" json:"duration"`

	Views int64 `bson:"views" json:"views"`

	Visibility Visibility `bson:"visibility" json:"visibility"`

	LikesCount int64 `bson:"likesCount" json:"likesCount"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type UploadURLRequest struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
}

type ConfirmUploadRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	S3Key       string `json:"s3Key"`
}

type CompleteVideoRequest struct {
	HLSURL       string `json:"hlsUrl"`
	ThumbnailKey string `json:"thumbnailKey"`
}

type OwnerResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VideoResponse struct {
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	HLSURL       string          `json:"hlsUrl"`
	ThumbnailKey string          `json:"thumbnailKey"`
	Status       VideoStatus     `json:"status"`
	Owner         *OwnerResponse `json:"owner,omitempty"`
}

type VideoListResponse struct {
	ID           string          `json:"_id"`
	Title        string          `json:"title"`
	ThumbnailKey string          `json:"thumbnailKey"`
	HLSURL       string          `json:"hlsUrl"`
	CreatedAt    time.Time       `json:"createdAt"`
	Views        int64           `json:"views"`
	Owner         *OwnerResponse `json:"owner,omitempty"`
}