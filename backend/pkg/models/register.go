package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	Age         int       `json:"age" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password_hash" validate:"required"`
	FirstName   string    `json:"first_name" validate:"required"`
	LastName    string    `json:"last_name" validate:"required"`
	Role        string    `json:"role" validate:"required"`
	Gender      string    `json:"gender" validate:"required"`
	DateOfBirth string    `json:"date_of_birth" validate:"required"`
	Avatar      string    `json:"avatar,omitempty"`
	Bio         string    `json:"bio,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Address     string    `json:"address,omitempty"`
	IsPrivate   bool      `json:"is_private" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
