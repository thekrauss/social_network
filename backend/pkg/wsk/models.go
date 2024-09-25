package wsk

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
