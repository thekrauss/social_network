package wsk

import (
	"backend/pkg/db"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upGradeWebsocket = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

func CheckOrigin(r *http.Request) bool {
	log.Printf("%s %s %s %v", r.Method, r.Host, r.RequestURI, r.Proto)
	return r.Method == http.MethodGet
}

func (w *WebsocketChat) UsersChatManager() {
	for {
		select {
		case user := <-w.JoinChannel:
			w.Mu.Lock()
			w.Users[user.Username] = user
			w.sendHistory(user)
			w.Mu.Unlock()

		case user := <-w.LeaveChannel:
			w.Mu.Lock()
			delete(w.Users, user.Username)
			w.Mu.Unlock()

		case msg := <-w.MessageChannel:
			w.Mu.Lock()

			if msg.MessageType == "typing" {
				// Gestion des événements "typing"
				w.handleTypingEvent(msg)
			} else if msg.MessageType == "message" && (msg.Content != "" || msg.Emoji != "") {
				// Gestion des messages réels
				w.handleMessageEvent(msg)
			}

			w.Mu.Unlock()
		}
	}
}

func (w *WebsocketChat) handleTypingEvent(msg *Message) {

	if targetUser, ok := w.Users[msg.RecipientID.String()]; ok {

		log.Printf("%s is typing...", msg.SenderID)
		err := targetUser.Connection.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending typing status to %s: %v", targetUser.Username, err)
		}

	}
}

func (w *WebsocketChat) handleMessageEvent(msg *Message) {

	// verifie si l'un des deux utilisateurs se suit ou si le profil est public
	if w.canSendMessage(msg.SenderID, msg.RecipientID) {
		w.sendPrivateMessage(msg)
		w.saveMessageHistory(msg)
	} else {
		log.Printf("Message blocked: %s cannot send message to %s", msg.SenderID, msg.RecipientID)
	}

}

func (w *WebsocketChat) canSendMessage(senderID, recipientID uuid.UUID) bool {

	recipient := w.Users[recipientID.String()]
	return recipient.IsPublic || w.areFollowingEachOther(senderID, recipientID)
}

func (w *WebsocketChat) areFollowingEachOther(userID1, userID2 uuid.UUID) bool {
	db, err := db.Store.OpenDatabase(&db.DBStore{})
	if err != nil {
		log.Println("Failed to open database in areFollowingEachOther:", err)
		return false
	}
	defer db.Close()

	var count int

	// Requête SQL pour vérifier si l'un des utilisateurs suit l'autre
	query := `
        SELECT COUNT(*) 
        FROM followers 
        WHERE (follower_id = ? AND followed_id = ? AND status = 'accepted') 
           OR (follower_id = ? AND followed_id = ? AND status = 'accepted')`

	err = db.QueryRow(query, userID1, userID2, userID2, userID1).Scan(&count)
	if err != nil {
		log.Println("Error querying the database in areFollowingEachOther:", err)
		return false
	}

	// count > 0, alors il y a une relation de suivi dans une direction
	return count > 0
}

func (w *WebsocketChat) sendPrivateMessage(msg *Message) {

	if targetUser, ok := w.Users[msg.RecipientID.String()]; ok {
		log.Printf("Sending message from %s to %s", msg.SenderID, msg.RecipientID)
		err := targetUser.Connection.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to %s: %v", targetUser.Username, err)
		}
	}
}

func (w *WebsocketChat) saveMessageHistory(msg *Message) {

	// Sauvegarder l'historique du message pour l'expéditeur et le destinataire
	w.MessageHistory[msg.SenderID.String()] = append(w.MessageHistory[msg.SenderID.String()], msg)
	w.MessageHistory[msg.RecipientID.String()] = append(w.MessageHistory[msg.RecipientID.String()], msg)
}

func (w *WebsocketChat) sendHistory(user *UserChat) {

	if messages, ok := w.MessageHistory[user.Username]; ok {
		for _, msg := range messages {
			user.Connection.WriteJSON(msg)
		}
	}
}
