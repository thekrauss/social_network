package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	PostID    uuid.UUID `json:"post_id" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	CreatedAt time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
}

type CommentPostGroup struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	PostID    uuid.UUID `json:"post_id" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	CreatedAt time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
}
