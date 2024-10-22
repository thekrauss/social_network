package controllers

import (
	"backend/pkg/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreatePostGroupHandler() http.HandlerFunc {
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

		var postGroup models.PostGroup
		if err := json.NewDecoder(r.Body).Decode(&postGroup); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Assigner un nouvel ID et d'autres champs à la publication
		postGroup.ID = uuid.Must(uuid.NewV4())
		postGroup.UserID = userID
		postGroup.CreatedAt = time.Now()
		postGroup.UpdatedAt = time.Now()

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Insérer la publication dans la base de données
		query := `INSERT INTO group_posts (id, group_id, user_id, title, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
		_, err = DB.Exec(query, postGroup.ID, postGroup.GroupID, postGroup.UserID, postGroup.Title, postGroup.Content, postGroup.CreatedAt, postGroup.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Post created successfully"))
	}
}

func (s *MyServer) ListPostGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'ID du groupe depuis les paramètres de l'URL
		groupIDStr := r.URL.Query().Get("group_id")
		if groupIDStr == "" {
			http.Error(w, "Group ID not provided", http.StatusBadRequest)
			return
		}

		groupID, err := uuid.FromString(groupIDStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
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
		// Fetch les posts depuis la base de données avec pagination
		log.Printf("Fetching posts from database (page: %d, limit: %d)\n", page, limit)

		query := `SELECT id, group_id, user_id, title, content, created_at, updated_at FROM group_posts WHERE group_id = ? LIMIT ? OFFSET ?`
		rows, err := DB.Query(query, groupID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var postsGroup []models.PostGroup
		for rows.Next() {
			var postgroup models.PostGroup
			if err := rows.Scan(&postgroup.ID, &postgroup.GroupID, &postgroup.UserID, &postgroup.Title, &postgroup.Content, &postgroup.CreatedAt, &postgroup.UpdatedAt); err != nil {
				http.Error(w, "Failed to scan post", http.StatusInternalServerError)
				return
			}
			postsGroup = append(postsGroup, postgroup)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(postsGroup); err != nil {
			http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		}
	}
}
