package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) LikeCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get comment ID from the form
		commentIDStr := r.FormValue("comment_id")
		if commentIDStr == "" {
			http.Error(w, "Comment ID not provided", http.StatusBadRequest)
			return
		}

		// Convert commentID to UUID
		commentID, err := uuid.FromString(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Get user ID from session or context (adjust to your session handling logic)
		userID, err := s.getUserIDFromSession(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Open database connection
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Begin a transaction
		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		// Toggle like for the comment
		err = ToggleLikeComment(userID, commentID, DB, tx)
		if err != nil {
			tx.Rollback() // Rollback transaction in case of error
			http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// Send success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment like toggled successfully"})
	}
}

func (s *MyServer) UnlikeCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse request form
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get comment ID from the form
		commentIDStr := r.FormValue("comment_id")
		if commentIDStr == "" {
			http.Error(w, "Comment ID not provided", http.StatusBadRequest)
			return
		}

		// Convert commentID to UUID
		commentID, err := uuid.FromString(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// Get user ID from session
		userID, err := s.getUserIDFromSession(r)
		if err != nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		// Open database connection
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// Begin a transaction
		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		// Toggle unlike for the comment
		err = ToggleUnLikeComment(userID, commentID, DB, tx)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to toggle unlike", http.StatusInternalServerError)
			return
		}

		// Commit transaction
		err = tx.Commit()
		if err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment unlike toggled successfully"})
	}
}

func (s *MyServer) getUserIDFromSession(r *http.Request) (uuid.UUID, error) {
	// Récupérer le cookie "user_id"
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return uuid.Nil, fmt.Errorf("no user ID cookie found")
	}

	// Convertir la valeur du cookie en UUID
	userID, err := uuid.FromString(cookie.Value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID cookie value")
	}

	return userID, nil
}

func ToggleLikeComment(userID, commentID uuid.UUID, DB *sql.DB, tx *sql.Tx) error {
	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	if liked {
		err := DeleteLikeComment(userID, commentID, tx)
		if err != nil {
			return err
		}
	} else {
		if unliked {
			err := DeleteUnLikeComment(userID, commentID, tx)
			if err != nil {
				return err
			}
		}
		err := CreateLikeComment(userID, commentID, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func ToggleUnLikeComment(userID, commentID uuid.UUID, DB *sql.DB, tx *sql.Tx) error {

	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user liked comment:", err)
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		log.Println("Error checking if user unliked comment:", err)
		return err
	}

	if unliked {
		err := DeleteUnLikeComment(userID, commentID, tx)
		if err != nil {
			log.Println("Error delete unlike comment:", err)
			return err
		}
	} else {
		if liked {
			err := DeleteLikeComment(userID, commentID, tx)
			if err != nil {
				log.Println("Error delete like comment:", err)
				return err
			}
		}
		err := CreateUnLikeComment(userID, commentID, tx)
		if err != nil {
			log.Println("Error create unlike:", err)
			return err
		}
	}

	return nil
}

/*--------------------------------------------------------------------*/

func UpdateCommentLikeCount(tx *sql.Tx, commentID uuid.UUID, likeType string, increment bool) error {
	operation := "+"
	if !increment {
		operation = "-"
	}

	query := fmt.Sprintf("UPDATE comments SET total_%s = total_%s %s 1 WHERE id = ?", likeType, likeType, operation)
	_, err := tx.Exec(query, commentID)
	if err != nil {
		log.Printf("Error updating total %s comments: %v", likeType, err)
		return err
	}

	return nil
}

func DeleteLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM commentlikes WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		log.Println("Error deleting like comment:", err)
		return err
	}

	return UpdateCommentLikeCount(tx, commentID, "likes", false)
}

func CreateLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO commentlikes (user_id, comment_id) VALUES (?, ?)", userID, commentID)
	if err != nil {
		log.Println("Error creating like comment:", err)
		return err
	}

	return UpdateCommentLikeCount(tx, commentID, "likes", true)
}

func DeleteUnLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM unlikescomment WHERE user_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		log.Println("Error deleting unlike comment:", err)
		return err
	}

	return UpdateCommentLikeCount(tx, commentID, "unlikescomment", false)
}

func CreateUnLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO unlikescomment (user_id, comment_id) VALUES (?, ?)", userID, commentID)
	if err != nil {
		log.Println("Error creating unlike comment:", err)
		return err
	}

	return UpdateCommentLikeCount(tx, commentID, "unlikescomment", true)
}

func UserUnLikedComment(userID, commentID uuid.UUID, tx *sql.Tx) (bool, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM unlikescomment WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func UserLikedComment(userID, commentID uuid.UUID, tx *sql.Tx) (bool, error) {
	// verifie si l'utilisateur a aimé le commentaire
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM commentlikes WHERE user_id = ? AND comment_id = ?", userID, commentID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
