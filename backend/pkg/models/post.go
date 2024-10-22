package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type Post struct {
	ID           uuid.UUID   `json:"id" validate:"required"`
	Title        string      `json:"title" validate:"required"`
	Category     string      `json:"category" validate:"required"`
	Content      string      `json:"content" validate:"required"`
	UserID       uuid.UUID   `json:"user_id" validate:"required"`
	Visibility   string      `json:"visibility" validate:"oneof=public private limited" default:"public"`
	CreatedAt    time.Time   `json:"created_at" default:"CURRENT_TIMESTAMP"`
	ImagePath    string      `json:"image_path,omitempty"`
	Username     string      `json:"username" validate:"required"`
	AllowedUsers []uuid.UUID `json:"allowed_users,omitempty"` // Utilisateurs autorisés pour les posts "almost_private"
}

type PostGroup struct {
	ID        uuid.UUID `json:"id"`
	GroupID   uuid.UUID `json:"group_id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
