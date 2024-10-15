package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
)

func GetPostsWithPagination(db *sql.DB, limit int, offset int) ([]models.Post, error) {
	query := `SELECT id, title, content, image_path, user_id, created_at FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts: %w", err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImagePath, &post.UserID, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}

func GetVisiblePostsWithPagination(db *sql.DB, userID uuid.UUID, limit int, offset int) ([]models.Post, error) {
	query := `
		SELECT p.id, p.title, p.content, p.image_path, p.visibility, p.created_at
		FROM posts p
		LEFT JOIN followers f ON p.user_id = f.followed_id AND f.follower_id = ?
		LEFT JOIN post_allowed_users pa ON p.id = pa.post_id AND pa.user_id = ?
		WHERE p.visibility = 'public' 
		OR (p.visibility = 'private' AND f.status = 'accepted')
		OR (p.visibility = 'almost_private' AND pa.user_id IS NOT NULL)
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.Query(query, userID, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImagePath, &post.Visibility, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *MyServer) StorePost(post models.Post) (uuid.UUID, error) {
	DB, err := s.Store.OpenDatabase()
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to open database: %v", err)
	}
	defer DB.Close()

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
