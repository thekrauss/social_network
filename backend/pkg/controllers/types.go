package controllers

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	DateOfBirth string    `json:"date_of_birth"` // Assuming string format for now
	Avatar      string    `json:"avatar"`        // URL or path to the avatar image
	Bio         string    `json:"bio"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Response struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
