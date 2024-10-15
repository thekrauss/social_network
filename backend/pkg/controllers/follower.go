package controllers

import (
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) FollowUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// récupérer l'utilisateur qui fait la demande et l'utilisateur cible
		senderID := r.Context().Value("userID").(uuid.UUID)
		receiverID, err := uuid.FromString(r.FormValue("receiver_id"))
		if err != nil {
			log.Println("Invalid UUID format for receiver_id:", err)
			http.Error(w, "Invalid receiver_id format", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// vérifier si l'utilisateur cible a un profil public ou privé
		var isPrivate bool
		err = DB.QueryRow("SELECT is_private FROM users WHERE id = ?", receiverID).Scan(&isPrivate)
		if err != nil {
			log.Println("Failed to retrieve user profile status", err)
			http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
			return
		}

		if !isPrivate {
			// si le profil est public ajouter directement à la table des followers
			_, err := DB.Exec("INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)", senderID, receiverID)
			if err != nil {
				log.Println("Failed to follow user", err)
				http.Error(w, "Failed to follow user", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("You are now following this user"))
		} else {
			// Si le profil est privé, ajouter une demande de suivi dans follow_requests
			_, err := DB.Exec("INSERT INTO follow_requests (sender_id, receiver_id) VALUES (?, ?)", senderID, receiverID)
			if err != nil {
				log.Println("Failed to send follow request", err)
				http.Error(w, "Failed to send follow request", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Follow request sent"))
		}
	}
}

func (s *MyServer) HandleFollowRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		receiverID := r.Context().Value("userID").(uuid.UUID)      // Utilisateur qui reçoit la demande
		senderID, err := uuid.FromString(r.FormValue("sender_id")) // Utilisateur qui a envoyé la demande
		if err != nil {
			log.Println("Invalid UUID format for sender_id:", err)
			http.Error(w, "Invalid sender_id format", http.StatusBadRequest)
			return
		}
		action := r.FormValue("action")

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		if action == "accept" {
			// Accepter la demande : supprimer la demande et ajouter l'entrée dans followers
			tx, err := DB.Begin()
			if err != nil {
				log.Println("Failed to begin transaction", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			_, err = tx.Exec("DELETE FROM follow_requests WHERE sender_id = ? AND receiver_id = ?", senderID, receiverID)
			if err != nil {
				tx.Rollback()
				log.Println("Failed to delete follow request", err)
				http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
				return
			}

			_, err = tx.Exec("INSERT INTO followers (follower_id, followed_id) VALUES (?, ?)", senderID, receiverID)
			if err != nil {
				tx.Rollback()
				log.Println("Failed to insert follower", err)
				http.Error(w, "Failed to accept follow request", http.StatusInternalServerError)
				return
			}

			tx.Commit()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Follow request accepted"))

		} else if action == "refuse" {

			// Refuser la demande : supprimer la demande dans follow_requests
			_, err := DB.Exec("DELETE FROM follow_requests WHERE sender_id = ? AND receiver_id = ?", senderID, receiverID)
			if err != nil {
				log.Println("Failed to delete follow request", err)
				http.Error(w, "Failed to decline follow request", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Follow request declined"))

		} else {
			http.Error(w, "Invalid action", http.StatusBadRequest)
		}
	}
}

func (s *MyServer) UnfollowUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		followerID := r.Context().Value("userID").(uuid.UUID)          // Utilisateur qui se désabonne
		followedID, err := uuid.FromString(r.FormValue("followed_id")) // Utilisateur à ne plus suivre
		if err != nil {
			log.Println("Invalid UUID format for followed_id:", err)
			http.Error(w, "Invalid followed_id format", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		_, err = DB.Exec("DELETE FROM followers WHERE follower_id = ? AND followed_id = ?", followerID, followedID)
		if err != nil {
			log.Println("Failed to unfollow user", err)
			http.Error(w, "Failed to unfollow user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unfollowed user successfully"))
	}
}
