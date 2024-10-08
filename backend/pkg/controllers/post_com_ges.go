package controllers

import (
	"backend/pkg/db"
	"backend/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

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

/*----------------------------------------------------------------------------------------------------------------*/

func GetCommentsByPost(DB *sql.DB, postID uuid.UUID) ([]models.Comment, error) {
	if DB == nil {
		return nil, errors.New("database connection is nil")
	}

	rows, err := DB.Query("SELECT id, content, post_id, user_id, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.PostID, &comment.UserID, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		//  le nom d'utilisateur associé
		username, err := GetUsernameByID(DB, comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.Username = username
		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *MyServer) StoreComment(comment models.Comment) error {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer DB.Close()

	query := `INSERT INTO comments (id, post_id, content, user_id, username, created_at)
              VALUES (?, ?, ?, ?, ?, ?)`
	_, err = DB.Exec(query, comment.ID, comment.PostID, comment.Content, comment.UserID, comment.Username, comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert comment into database: %v", err)
	}

	return nil
}

func (s *MyServer) PostExists(postID uuid.UUID) bool {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		log.Println("Failed to open database:", err)
		return false
	}
	defer DB.Close()

	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check if post exists for post ID %s: %v", postID, err)
		return false
	}
	return exists
}
