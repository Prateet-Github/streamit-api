package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Comment struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	UserID  bson.ObjectID `bson:"userId" json:"userId"`
	VideoID bson.ObjectID `bson:"videoId" json:"videoId"`

	Content string `bson:"content" json:"content"`

	ParentCommentID *bson.ObjectID `bson:"parentCommentId,omitempty" json:"parentCommentId,omitempty"`

	RepliesCount int64 `bson:"repliesCount" json:"repliesCount"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type CreateCommentRequest struct {
	Content         string `json:"content" binding:"required,max=500"`
	ParentCommentID string `json:"parentCommentId,omitempty"`
}

type ReplyCommentRequest struct {
	Content string `json:"content" binding:"required,max=500"`
}

type CommentResponse struct {
	ID string `json:"id"`

	Content string `json:"content"`

	RepliesCount int64 `json:"repliesCount"`

	CreatedAt time.Time `json:"createdAt"`

	User *OwnerResponse `json:"user"`
}
