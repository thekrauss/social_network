package controllers

import (
	"backend/pkg/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreateCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var comment models.Comment

			if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
				log.Printf("Failed to decode comment request payload: %v", err)
				http.Error(w, "Invalid request comment payload", http.StatusBadRequest)
				return
			}

			if comment.PostID == uuid.Nil {
				log.Println("Invalid post ID")
				http.Error(w, "Invalid post ID", http.StatusBadRequest)
				return
			}

			// Vérifier si le post_id existe dans la base de données
			if !s.PostExists(comment.PostID) {
				log.Println("Post ID does not exist")
				http.Error(w, "Post ID does not exist", http.StatusBadRequest)
				return
			}
			comment.CreatedAt = time.Now()
			userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
			if !ok || userID == uuid.Nil {
				log.Println("User ID not found in context")
				http.Error(w, "User ID not found in context", http.StatusUnauthorized)
				return
			}

			comment.UserID = userID

			fmt.Printf("userID: %v\n", comment.UserID)

			if err := s.StoreComment(comment); err != nil {
				log.Println("Failed to store comment:", err)
				http.Error(w, "Failed to store comment", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(comment)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (s *MyServer) ListCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//  l'ID du post depuis l'URL
		query := r.URL.Query()
		postIDStr := query.Get("post_id")
		if postIDStr == "" {
			http.Error(w, "Post ID not provided", http.StatusBadRequest)
			return
		}

		//  postID en UUID
		postID, err := uuid.FromString(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		//  les commentaires liés au postID
		comments, err := GetCommentsByPost(DB, postID)
		if err != nil {
			http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"comments": comments,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode comments to JSON:", err)
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		}
	}
}
