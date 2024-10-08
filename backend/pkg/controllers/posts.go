package controllers

import (
	"backend/pkg/db"
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) PostHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var post models.Post

		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			log.Println("Failed to decode post request payload", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if len(post.Title) == 0 || len(post.Content) == 0 {
			log.Println("Title or content empty")
			http.Error(w, "Title or content empty", http.StatusBadRequest)
			return
		}

		// gestion du téléchargement d'image
		var imagesPath string
		file, handler, err := r.FormFile("image")
		if err == nil { // S'il y a un fichier image dans la requête
			defer file.Close()

			// vérification de l'extension de l'image
			if !IsValidImageExtension(handler.Filename) {
				log.Println("Invalid image file extension")
				http.Error(w, "Invalid image file extension", http.StatusBadRequest)
				return
			}

			imagesPath, err = UploadImages(w, r, "./image_path/")
			if err != nil {
				log.Println("Failed to upload image:", err)
				http.Error(w, "Failed to upload image", http.StatusInternalServerError)
				return
			}
			post.ImagePath = imagesPath
		}

		//  l'ID de l'utilisateur à partir du contexte en tant qu'UUID
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

		// envoi  la réponse de succès
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

func (s *MyServer) StorePost(post models.Post) (uuid.UUID, error) {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer DB.Close()
	//  si la table 'posts' existe ou la créer
	_, err = DB.Exec(db.Posts_table)
	if err != nil {
		log.Println("Error ensuring posts table exists:", err)
		return uuid.Nil, fmt.Errorf("failed to ensure posts table exists: %v", err)
	}
	log.Println("Database and table ready")

	//  l'UUID pour le nouveau post
	postID := uuid.Must(uuid.NewV4())
	query := `INSERT INTO posts (id, user_id, title, content, image_path)
	VALUES (?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, postID, post.UserID, post.Title, post.Content, post.ImagePath)
	if err != nil {
		log.Println("Failed to insert post into database:", err)
		return uuid.Nil, fmt.Errorf("failed to insert post: %v", err)
	}

	log.Println("Post successfully created with ID:", postID)
	return postID, nil
}

func (s *MyServer) ListPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("failed to open database for ListPost", err)
			http.Error(w, "failed to open database for ListPost", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		log.Println("Fetching posts from database")
		posts, err := GetPosts(DB)
		if err != nil {
			log.Println("Failed to retrieve posts:", err)
			http.Error(w, "Failed to retrieve posts from the database", http.StatusInternalServerError)
			return
		}
		log.Printf("Retrieved %d posts from the database\n", len(posts))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("Failed to encode posts to JSON:", err)
			http.Error(w, "Failed to encode posts to JSON", http.StatusInternalServerError)
		}
	}
}

func GetPosts(DB *sql.DB) ([]models.Post, error) {
	if DB == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := DB.Query("SELECT p.id, p.title, p.content, p.image_path, p.user_id, p.created_at, p.category, u.username FROM posts p JOIN users u ON p.user_id = u.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImagePath, &post.UserID, &post.CreatedAt, &post.Category, &post.Username)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if len(posts) == 0 {
		log.Println("No posts found")
	}
	return posts, nil
}
