package models

import (
	"github.com/gofrs/uuid"
)

// structure de base d'un groupe
type Group struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description" validate:"required"`
	CreatorID   uuid.UUID     `json:"creator_id"`
	Members     []GroupMember `json:"members,omitempty"` // liste des membres du groupe
	CreatedAt   string        `json:"created_at"`
}

// structure pour les membres du groupe
type GroupMember struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role" validate:"oneof=creator member"`
	Status string    `json:"status" validate:"oneof=pending accepted"`
}
