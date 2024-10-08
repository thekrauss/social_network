package wsk

/*
import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID             int       `json:"id"`
	SenderUsername string    `json:"senderUsername"`
	TargetUsername string    `json:"targetUsername"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	IsTyping       bool      `json:"isTyping"`
	Type           string    `json:"type"`
}

type Message struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	SenderID   uuid.UUID `json:"sender_id" validate:"required"`
	RecipientID uuid.UUID `json:"recipient_id" validate:"required"`
	Content    string    `json:"content" validate:"required"`
	CreatedAt  time.Time `json:"created_at" default:"CURRENT_TIMESTAMP"`
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
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type UserChat struct {
	Channels   *Channel
	Username   string
	Connection *websocket.Conn
}

type userChannel chan *userChannel
type messageChannel chan *messageChannel

type Channel struct {
	messageChannel messageChannel
	leaveChannel   userChannel
}

/*----------------------------------------------------------*/
/*
func NewUserChat(channels *Channel, username string, conn *websocket.Conn) *UserChat {
	return &UserChat{
		Channels:   channels,
		Username:   username,
		Connection: conn,
	}
}

func NewWebsocketChat() *WebsocketChat {
	w := &WebsocketChat{
		Users:          make(map[string]*UserChat),
		JoinChannel:    make(userChannel),
		LeaveChannel:   make(userChannel),
		MessageChannel: make(messageChannel),
		MessageHistory: make(map[string][]*Message),
	}
	go w.UsersChatManager()
	return w
}
*/
