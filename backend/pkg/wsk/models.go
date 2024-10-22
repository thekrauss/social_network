package wsk

import (
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	ID             uuid.UUID `json:"id" validate:"required"`
	SenderID       uuid.UUID `json:"sender_id" validate:"required"`
	SenderUsername string    `json:"sender_username" validate:"required"`
	RecipientID    uuid.UUID `json:"recipient_id,omitempty"`
	GroupID        uuid.UUID `json:"group_id,omitempty"`
	Content        string    `json:"content" validate:"required"`
	Emoji          string    `json:"emoji,omitempty"`
	CreatedAt      time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
	Timestamp      time.Time `json:"timestamp"`
	MessageType    string    `json:"message_type" validate:"oneof=text emoji"`
}

type Notification struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
	Read      bool      `json:"read" default:"false"`
}

type Follower struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	FollowerID uuid.UUID `json:"follower_id" validate:"required"`
	FollowedID uuid.UUID `json:"followed_id" validate:"required"`
	Status     string    `json:"status" validate:"oneof=pending accepted" default:"pending"`
	CreatedAt  time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
}

type Group struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description,omitempty"`
	CreatorID   uuid.UUID `json:"creator_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
}

type GroupMember struct {
	ID        uuid.UUID `json:"id" validate:"required"`
	GroupID   uuid.UUID `json:"group_id" validate:"required"`
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	Status    string    `json:"status" validate:"oneof=pending accepted" default:"pending"`
	CreatedAt time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
}

type WebsocketChat struct {
	Users          map[string]*UserChat
	JoinChannel    userChannel
	LeaveChannel   userChannel
	MessageChannel messageChannel
	MessageHistory map[string][]*Message
	Mu             sync.Mutex
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	IsPublic bool      `json:"is_public"`
}

type UserChat struct {
	Channels   *Channel
	Username   string
	Connection *websocket.Conn
	IsPublic   bool `json:"is_public"`
}

type Channel struct {
	MessageChannel messageChannel
	LeaveChannel   userChannel
}

type userChannel chan *UserChat
type messageChannel chan *Message
