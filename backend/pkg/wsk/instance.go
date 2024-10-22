package wsk

import "github.com/gorilla/websocket"

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

func NewUserChat(channels *Channel, username string, conn *websocket.Conn) *UserChat {
	return &UserChat{
		Channels:   channels,
		Username:   username,
		Connection: conn,
	}
}
