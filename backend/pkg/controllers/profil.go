package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func (s *MyServer) MyProfil() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			log.Println("User ID not found in context")
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database for MyProfil", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Commencer une transaction
		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		commitErr := func() error {
			if err != nil {
				return tx.Rollback()
			}
			return tx.Commit()
		}

		defer func() {
			if commitErr() != nil {
				log.Println("Transaction failed to commit/rollback")
			}
		}()

		// recupere les paramètres de pagination
		queryParams := r.URL.Query()
		limit, err := strconv.Atoi(queryParams.Get("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(queryParams.Get("offset"))
		if err != nil || offset < 0 {
			offset = 0
		}

		// recupere le profil utilisateur avec pagination pour les posts
		response, err := GetMyProfil(DB, userID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to get MyProfil", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode response to JSON:", err)
			http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		}
	}
}

func GetMyProfil(db *sql.DB, userID uuid.UUID, limit int, offset int) (models.UserProfil, error) {
	var profil models.UserProfil

	query := `SELECT id, username, first_name, last_name, bio FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&profil.UserID, &profil.Username, &profil.FirstName, &profil.LastName, &profil.Bio)
	if err != nil {
		return profil, fmt.Errorf("failed to query user Profil: %w", err)
	}

	// Récupérer les followers
	profil.Followers, err = GetFollowers(db, userID)
	if err != nil {
		return profil, fmt.Errorf("failed to get followers: %w", err)
	}

	// Récupérer les utilisateurs suivis
	profil.Following, err = GetFollowing(db, userID)
	if err != nil {
		return profil, fmt.Errorf("failed to get following: %w", err)
	}

	// Récupérer les posts de l'utilisateur avec pagination
	profil.Posts, err = GetProfilPostsWithPagination(db, limit, offset)
	if err != nil {
		return profil, fmt.Errorf("failed to get user posts: %w", err)
	}

	return profil, nil
}
