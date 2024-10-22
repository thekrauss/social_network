package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) LikeComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil || r.FormValue("comment_id") == "" {
			http.Error(w, "Comment ID not provided", http.StatusBadRequest)
			return
		}
		commentID, err := uuid.FromString(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		if err := ToggleLikeComment(userID, commentID, tx); err != nil {
			tx.Rollback()
			http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment like toggled successfully"})
	}
}

func (s *MyServer) UnLikeComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil || r.FormValue("comment_id") == "" {
			http.Error(w, "Comment ID not provided", http.StatusBadRequest)
			return
		}
		commentID, err := uuid.FromString(r.FormValue("comment_id"))
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}

		if err := ToggleUnLikeComment(userID, commentID, tx); err != nil {
			tx.Rollback()
			http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment like toggled successfully"})
	}
}

func ToggleLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	if liked {
		return DeleteLikeComment(userID, commentID, tx)
	}

	if unliked {
		if err := DeleteUnLikeComment(userID, commentID, tx); err != nil {
			return err
		}
	}

	return CreateLikeComment(userID, commentID, tx)
}

func ToggleUnLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	liked, err := UserLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	unliked, err := UserUnLikedComment(userID, commentID, tx)
	if err != nil {
		return err
	}

	if liked {
		return DeleteLikeComment(userID, commentID, tx)
	}

	if unliked {
		return DeleteUnLikeComment(userID, commentID, tx)
	}

	return CreateUnLikeComment(userID, commentID, tx)
}

/*--------------------------------------------------------------------*/

func CreateLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO comment_interactions (user_id, comment_id, interaction_type) VALUES (?, ?, 'like')", userID, commentID)
	if err != nil {
		return err
	}
	return UpdateCommentLikeCount(tx, commentID, "likes", true)
}

func DeleteLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM comment_interactions WHERE user_id = ? AND comment_id = ? AND interaction_type = 'like'", userID, commentID)
	if err != nil {
		return err
	}
	return UpdateCommentLikeCount(tx, commentID, "likes", false)
}

func CreateUnLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO comment_interactions (user_id, comment_id, interaction_type) VALUES (?, ?, 'unlike')", userID, commentID)
	if err != nil {
		return err
	}
	return UpdateCommentLikeCount(tx, commentID, "unlikes", true)
}

func DeleteUnLikeComment(userID, commentID uuid.UUID, tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM comment_interactions WHERE user_id = ? AND comment_id = ? AND interaction_type = 'unlike'", userID, commentID)
	if err != nil {
		return err
	}
	return UpdateCommentLikeCount(tx, commentID, "unlikes", false)
}

func UpdateCommentLikeCount(tx *sql.Tx, commentID uuid.UUID, likeType string, increment bool) error {
	operation := "+"
	if !increment {
		operation = "-"
	}

	query := fmt.Sprintf("UPDATE comments SET total_%s = total_%s %s 1 WHERE id = ?", likeType, likeType, operation)
	_, err := tx.Exec(query, commentID)
	return err
}

/*--------------------------------------------------------------------*/

func UserLikedComment(userID, commentID uuid.UUID, tx *sql.Tx) (bool, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE user_id = ? AND comment_id = ? AND interaction_type = 'like'", userID, commentID).Scan(&count)
	return count > 0, err
}

func UserUnLikedComment(userID, commentID uuid.UUID, tx *sql.Tx) (bool, error) {
	var count int
	err := tx.QueryRow("SELECT COUNT(*) FROM comment_interactions WHERE user_id = ? AND comment_id = ? AND interaction_type = 'unlike'", userID, commentID).Scan(&count)
	return count > 0, err
}
