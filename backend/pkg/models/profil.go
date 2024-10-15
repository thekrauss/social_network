package models

import "github.com/gofrs/uuid"

type UserProfil struct {
	UserID    uuid.UUID    `json:"user_id"`
	Username  string       `json:"username"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Bio       string       `json:"bio"`
	IsPrivate bool         `json:"is_private"`
	Followers []SimpleUser `json:"followers,omitempty"`
	Following []SimpleUser `json:"following,omitempty"`
	Posts     []Post       `json:"posts,omitempty"`
}

type SimpleUser struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
