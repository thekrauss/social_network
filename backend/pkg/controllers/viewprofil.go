package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) GetUserProfil() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.FromString(r.URL.Query().Get("user_id"))
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// récupérer l'utilisateur connecté pour déterminer si on peut afficher les données d'un profil privé
		loggedInUserID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database", err)
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

		log.Println("Database and table ready")

		userProfil, err := GetUserProfilFromDB(DB, userID, loggedInUserID)
		if err != nil {
			log.Println("Failed to get user profile", err)
			http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
			return
		}

		// Encoder la réponse en JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userProfil); err != nil {
			log.Println("Failed to encode user profile", err)
			http.Error(w, "Failed to encode user profile", http.StatusInternalServerError)
		}
	}
}

func GetUserProfilFromDB(db *sql.DB, userID, loggedInUserID uuid.UUID) (models.UserProfil, error) {

	var profil models.UserProfil

	query := `SELECT id, username, first_name, last_name, bio, is_private FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(&profil.UserID, &profil.Username, &profil.FirstName, &profil.LastName, &profil.Bio, &profil.IsPrivate)
	if err != nil {
		return profil, fmt.Errorf("failed to query user Profil: %w", err)
	}

	// si le profil est privé et l'utilisateur connecté n'est pas un follower, renvoyer une erreur
	if profil.IsPrivate && !IsUserFollower(db, userID, loggedInUserID) {
		return profil, nil
	}

	// récupére les followers
	profil.Followers, err = GetFollowers(db, userID)
	if err != nil {
		return profil, fmt.Errorf("failed to get followers: %w", err)
	}

	// récupére les utilisateurs suivis
	profil.Following, err = GetFollowing(db, userID)
	if err != nil {
		return profil, fmt.Errorf("failed to get following: %w", err)
	}

	profil.Posts, err = GetUserPosts(db, userID)
	if err != nil {
		return profil, fmt.Errorf("failed to get user posts: %w", err)
	}
	return profil, nil
}

func IsUserFollower(db *sql.DB, userID, followerID uuid.UUID) bool {
	var count int
	query := `SELECT COUNT(*) FROM followers WHERE followed_id = ? AND follower_id = ? AND status = 'accepted'`
	err := db.QueryRow(query, userID, followerID).Scan(&count)
	return err == nil && count > 0
}

func GetFollowers(db *sql.DB, userID uuid.UUID) ([]models.SimpleUser, error) {
	var followers []models.SimpleUser
	query := `SELECT u.id, u.username FROM users u INNER JOIN followers f ON u.id = f.follower_id WHERE f.followed_id = ? AND f.status = 'accepted'`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.SimpleUser
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		followers = append(followers, user)
	}
	return followers, nil
}

func GetFollowing(db *sql.DB, userID uuid.UUID) ([]models.SimpleUser, error) {
	var following []models.SimpleUser
	query := `SELECT u.id, u.username FROM users u INNER JOIN followers f ON u.id = f.followed_id WHERE f.follower_id = ? AND f.status = 'accepted'`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.SimpleUser
		if err := rows.Scan(&user.UserID, &user.Username); err != nil {
			return nil, err
		}
		following = append(following, user)
	}
	return following, nil
}

func GetUserPosts(db *sql.DB, userID uuid.UUID) ([]models.Post, error) {
	var posts []models.Post
	query := `SELECT id, title, content, created_at, visibility FROM posts WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.Visibility); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
