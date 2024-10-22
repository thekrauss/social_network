package wsk

import (
	"log"
	"time"

	"github.com/gofrs/uuid"
)

func (w *WebsocketChat) SendNotification(notification *Notification) {
	if targetUser, ok := w.Users[notification.UserID.String()]; ok {
		err := targetUser.Connection.WriteJSON(notification)
		if err != nil {
			log.Printf("Error sending notification to %s: %v", targetUser.Username, err)
		}
	}
}

func CreateNotification(userID uuid.UUID, content string) *Notification {
	return &Notification{
		ID:        uuid.Must(uuid.NewV4()),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		Read:      false,
	}
}

func (w *WebsocketChat) NotifyFollowRequest(followerID, followedID uuid.UUID) {
	content := "You have a new follow request!"
	notification := CreateNotification(followedID, content)
	w.SendNotification(notification)
}

func (w *WebsocketChat) NotifyGroupInvite(groupID, inviteeID uuid.UUID) {
	content := "You have been invited to join a group!"
	notification := CreateNotification(inviteeID, content)
	w.SendNotification(notification)
}

/*
func (w *WebsocketChat) NotifyGroupJoinRequest(groupID, requesterID uuid.UUID) {
	groupCreatorID := w.getGroupCreator(groupID) // Fonction fictive pour obtenir le créateur du groupe
	content := "A user has requested to join your group!"
	notification := CreateNotification(groupCreatorID, content)
	w.SendNotification(notification)
}
*/

func (w *WebsocketChat) NotifyGroupEvent(groupID, memberID uuid.UUID) {
	content := "A new event has been created in your group!"
	notification := CreateNotification(memberID, content)
	w.SendNotification(notification)
}
