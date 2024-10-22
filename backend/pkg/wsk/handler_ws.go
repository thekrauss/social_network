package wsk

import (
	"log"
	"net/http"
	"time"
)

func (w *WebsocketChat) HanderUsersConnection(wr http.ResponseWriter, r *http.Request) {
	webSocketConn, err := upGradeWebsocket.Upgrade(wr, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}

	cookie, err := r.Cookie("username")
	if err != nil {
		log.Printf("Failed to get username cookie HanderUsersConnection: %v", err)
		return
	}

	username := cookie.Value

	log.Printf("User %s connected", username)

	// creation d'une nouvelle session de chat utilisateur
	userChat := NewUserChat(&Channel{MessageChannel: w.MessageChannel, LeaveChannel: w.LeaveChannel}, username, webSocketConn)

	w.JoinChannel <- userChat

	// ecoute des messages entrants depuis cet utilisateur
	go userChat.listenForMessages()
}

func (u *UserChat) listenForMessages() {

	defer func() {
		// l'utilisateur se déconnecte, il est retiré du canal LeaveChannel
		u.Channels.LeaveChannel <- u
		u.Connection.Close()
		log.Printf("User %s disconnected", u.Username)
	}()

	for {

		var msg Message

		// lire les messages JSON envoyés par l'utilisateur via WebSocket
		err := u.Connection.ReadJSON(&msg)
		if err != nil {
			// si l'erreur est liée à une fermeture de connexion on arrête
			log.Printf("Error reading JSON from user %s: %v", u.Username, err)
			break
		}

		// assigner l'username de l'expéditeur et ajouter un timestamp au message
		msg.SenderUsername = u.Username
		msg.Timestamp = time.Now()

		// gestion typing ou message normal
		switch msg.MessageType {
		case "typing":
			log.Printf("%s is typing...", msg.SenderUsername)
			u.Channels.MessageChannel <- &msg
		case "message":
			// Ignorer les messages vides
			if msg.Content == "" && msg.Emoji == "" {
				log.Printf("Empty message from %s ignored", msg.SenderUsername)
				continue
			}
			log.Printf("Message to send: %+v", msg)
			u.Channels.MessageChannel <- &msg
		default:
			log.Printf("Unknown message type from %s: %s", u.Username, msg.MessageType)
		}

	}
}
