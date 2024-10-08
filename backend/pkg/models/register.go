package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" validate:"required"`            // UUID unique pour l'utilisateur
	Username    string    `json:"username" validate:"required"`      // Nom d'utilisateur obligatoire
	Age         int       `json:"age" validate:"required"`           // Âge obligatoire
	Gender      string    `json:"gender" validate:"required"`        // Genre obligatoire
	FirstName   string    `json:"firstname" validate:"required"`     // Prénom obligatoire
	LastName    string    `json:"lastname" validate:"required"`      // Nom de famille obligatoire
	Email       string    `json:"email" validate:"required,email"`   // Email obligatoire et validation de format email
	Password    string    `json:"password" validate:"required"`      // Mot de passe obligatoire
	Role        string    `json:"role" validate:"required"`          // Rôle obligatoire (par exemple admin, user)
	DateOfBirth string    `json:"date_of_birth" validate:"required"` // Date de naissance obligatoire
	Avatar      string    `json:"avatar,omitempty"`                  // Avatar optionnel
	Bio         string    `json:"bio,omitempty"`                     // Bio optionnelle
	Phone       string    `json:"phone,omitempty"`                   // Téléphone optionnel
	Address     string    `json:"address,omitempty"`                 // Adresse optionnelle
	IsPrivate   bool      `json:"is_private" validate:"required"`    // Visibilité du profil (privé ou public) obligatoire
	CreatedAt   time.Time `json:"created_at" validate:"required"`    // Date de création obligatoire
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`    // Date de dernière mise à jour obligatoire
}

type Response struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}
