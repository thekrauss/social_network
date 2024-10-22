package controllers

import (
	"backend/pkg/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreatePostHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(20 << 20) // Limite de 20MB pour les fichiers
		if err != nil {
			log.Println("Failed to parse multipart form:", err)
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		var post models.Post
		// décoder les informations du post depuis le JSON contenu dans le champ "data"
		if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&post); err != nil {
			log.Println("Failed to decode post request payload", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if len(post.Title) == 0 || len(post.Content) == 0 {
			log.Println("Title or content empty")
			http.Error(w, "Title or content empty", http.StatusBadRequest)
			return
		}

		if post.Visibility != "public" && post.Visibility != "private" && post.Visibility != "almost_private" {
			log.Println("Invalid post visibility")
			http.Error(w, "Invalid post visibility", http.StatusBadRequest)
			return
		}

		if post.Visibility == "almost_private" {
			allowedUsersStr := r.FormValue("allowed_users")
			if allowedUsersStr != "" {
				allowedUsers := strings.Split(allowedUsersStr, ",")
				for _, userIDStr := range allowedUsers {
					allowedUserID, err := uuid.FromString(userIDStr)
					if err != nil {
						log.Println("Invalid allowed user ID:", userIDStr)
						http.Error(w, "Invalid allowed user ID", http.StatusBadRequest)
						return
					}
					post.AllowedUsers = append(post.AllowedUsers, allowedUserID)
				}

			}
		}

		// gestion du téléchargement d'image
		file, handler, err := r.FormFile("image")
		if err == nil {
			defer file.Close()

			if !IsValidImageExtension(handler.Filename) {
				log.Println("Invalid image file extension")
				http.Error(w, "Invalid image file extension", http.StatusBadRequest)
				return
			}

			// Upload de l'image
			imagesPath, err := UploadImages(w, r, "./image_path/")
			if err != nil {
				log.Println("Failed to upload image:", err)
				http.Error(w, "Failed to upload image", http.StatusInternalServerError)
				return
			}
			post.ImagePath = imagesPath
		}

		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			log.Println("User ID not found in context")
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		post.UserID = userID

		postID, err := s.StorePost(post)
		if err != nil {
			log.Println("Failed to save post:", err)
			http.Error(w, "Failed to save post", http.StatusInternalServerError)
			return
		}

		post.ID = postID
		post.CreatedAt = time.Now()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func (s *MyServer) ListPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'ID utilisateur à partir du contexte
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			log.Println("User ID not found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("User ID found:", userID)

		// Ouvrir la base de données
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database for ListPost:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

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

		// Gestion des paramètres de pagination : page et limit
		page := 1
		limit := 10

		// Analyse des paramètres page et limit
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

		// Calculer l'offset pour la pagination
		offset := (page - 1) * limit

		log.Printf("Fetching visible posts from database (page: %d, limit: %d)\n", page, limit)
		posts, err := GetVisiblePostsWithPagination(DB, userID, limit, offset)
		if err != nil {
			log.Println("Failed to retrieve posts:", err)
			http.Error(w, "Failed to retrieve posts from the database", http.StatusInternalServerError)
			return
		}

		// Répondre avec la liste des posts
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("Failed to encode posts to JSON:", err)
			http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError)
		}
	}
}

/*--------------------------------------------------------------------------------------------------------------------------*/
