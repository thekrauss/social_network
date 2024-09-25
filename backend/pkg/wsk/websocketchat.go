package wsk

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
