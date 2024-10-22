package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// Structure de base d'un événement
type GroupEvent struct {
	ID          uuid.UUID    `json:"id"`
	GroupID     uuid.UUID    `json:"group_id"`
	UserID      uuid.UUID    `json:"user_id"`
	Title       string       `json:"title" validate:"required"`
	Description string       `json:"description" validate:"required"`
	EventDate   time.Time    `json:"event_date" validate:"required"` // Date et heure de l'événement
	Options     EventOptions `json:"options"`                        // Options de réponse à l'événement
	CreatedAt   time.Time    `json:"created_at"`
}

// Options de réponse pour l'événement
type EventOptions struct {
	Going    int `json:"going"`     // Nombre de participants qui ont répondu "Going"
	NotGoing int `json:"not_going"` // Nombre de participants qui ont répondu "Not going"
}

// Structure pour les réponses à l'événement
type EventResponse struct {
	EventID  uuid.UUID `json:"event_id"`
	UserID   uuid.UUID `json:"user_id"`
	Response string    `json:"response" validate:"oneof=Going NotGoing"`
}
