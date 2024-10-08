package models

import (
	"github.com/gofrs/uuid"
)

type CommentLike struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	CommentID uuid.UUID `json:"comment_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
}
