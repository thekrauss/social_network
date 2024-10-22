package controllers

import (
	"backend/pkg/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreateCommentPostsGroup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'utilisateur connecté
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		var comment models.CommentPostGroup
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Assigner un nouvel ID, l'ID de l'utilisateur et d'autres champs
		comment.ID = uuid.Must(uuid.NewV4())
		comment.UserID = userID
		comment.CreatedAt = time.Now()

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Insérer le commentaire dans la base de données
		query := `INSERT INTO group_posts_comments (id, post_id, content, user_id, username, created_at) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = DB.Exec(query, comment.ID, comment.PostID, comment.Content, comment.UserID, comment.Username, comment.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to create comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Comment created successfully"))
	}
}

func (s *MyServer) ListCommentsByPostGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'ID de la publication depuis les paramètres de l'URL
		postIDStr := r.URL.Query().Get("post_id")
		if postIDStr == "" {
			http.Error(w, "Post ID not provided", http.StatusBadRequest)
			return
		}

		postID, err := uuid.FromString(postIDStr)
		if err != nil {
			http.Error(w, "Invalid Post ID", http.StatusBadRequest)
			return
		}

		// Ouvrir la base de données
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// récupération des paramètres de pagination : page et limit
		page := 1
		limit := 10

		// Parse des paramètres page et limit
		queryParams := r.URL.Query()
		if p := queryParams.Get("page"); p != "" {
			page, err = strconv.Atoi(p)
			if err != nil || page < 1 {
				page = 1
			}
		}

		if l := queryParams.Get("limit"); l != "" {
			limit, err = strconv.Atoi(l)
			if err != nil || limit < 1 {
				limit = 10
			}
		}

		offset := (page - 1) * limit

		// Fetch les commentaires depuis la base de données avec pagination
		query := `SELECT id, post_id, content, user_id, username, created_at FROM group_posts_comments WHERE post_id = ? LIMIT ? OFFSET ?`
		rows, err := DB.Query(query, postID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var comments []models.CommentPostGroup
		for rows.Next() {
			var comment models.CommentPostGroup
			if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.UserID, &comment.Username, &comment.CreatedAt); err != nil {
				http.Error(w, "Failed to scan comment", http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}

		// Répondre avec la liste des commentaires
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			http.Error(w, "Failed to encode comments", http.StatusInternalServerError)
		}
	}
}
