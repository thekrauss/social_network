package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofrs/uuid"
)

func (s *MyServer) UpdateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		userID := r.Context().Value("userID").(uuid.UUID)

		var updatedUser models.User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			log.Println("Failed to decode request UpdateProfile:", err)
			http.Error(w, "Failed to decode request", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database UpdateProfile:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Nettoyer les champs texte
		updatedUser.FirstName = strings.TrimSpace(updatedUser.FirstName)
		updatedUser.LastName = strings.TrimSpace(updatedUser.LastName)
		updatedUser.Email = strings.TrimSpace(updatedUser.Email)

		if updatedUser.FirstName != "" && (len(updatedUser.FirstName) < 2 || len(updatedUser.FirstName) > 30) {
			http.Error(w, "First name must be between 2 and 30 characters", http.StatusBadRequest)
			return
		}

		if updatedUser.LastName != "" && (len(updatedUser.LastName) < 2 || len(updatedUser.LastName) > 30) {
			http.Error(w, "Last name must be between 2 and 30 characters", http.StatusBadRequest)
			return
		}

		if updatedUser.Gender != "" && !IsValidGender(updatedUser.Gender) {
			http.Error(w, "Invalid gender. Must be 'Homme' or 'Femme'.", http.StatusBadRequest)
			return
		}

		if updatedUser.Phone != "" && !IsValidPhoneNumber(updatedUser.Phone) {
			http.Error(w, "Invalid phone number format", http.StatusBadRequest)
			return
		}

		if updatedUser.Avatar != "" && !IsValidImageExtension(updatedUser.Avatar) {
			http.Error(w, "Invalid avatar format", http.StatusBadRequest)
			return
		}

		if updatedUser.Bio != "" && len(updatedUser.Bio) > 500 {
			http.Error(w, "Bio must be less than 500 characters", http.StatusBadRequest)
			return
		}

		if updatedUser.Address != "" && len(updatedUser.Address) > 100 {
			http.Error(w, "Address must be less than 100 characters", http.StatusBadRequest)
			return
		}

		// Vérification de l'utilisateur sans vérifier son propre email ou nom d'utilisateur
		err = checkUser(DB, updatedUser, userID)
		if err != nil {
			log.Println("Failed to check user:", err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		// Mise à jour des champs dans la base de données
		query := `UPDATE users SET first_name = ?, last_name = ?, email = ?, gender = ?, avatar = ?, bio = ?, phone_number = ?, address = ?, is_private = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = DB.Exec(query, updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.Gender, updatedUser.Avatar, updatedUser.Bio, updatedUser.Phone, updatedUser.Address, updatedUser.IsPrivate, userID)

		if err != nil {
			log.Println("Failed to update user profile:", err)
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		response := models.Response{Message: "Profile updated successfully"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func (s *MyServer) UpdateVisibility() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// récupére l'utilisateur connecté
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusUnauthorized)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// récupére le paramètre is_private depuis le corps de la requête
		var requestData struct {
			IsPrivate bool `json:"is_private"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// mettre à jour la visibilité du profil
		query := `UPDATE users SET is_private = ? WHERE id = ?`
		_, err = DB.Exec(query, requestData.IsPrivate, userID)
		if err != nil {
			log.Println("Failed to update profile visibility", err)
			http.Error(w, "Failed to update profile visibility", http.StatusInternalServerError)
			return
		}

		// Réponse
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profile visibility updated successfully"))
	}
}

func checkUser(db *sql.DB, user models.User, userID uuid.UUID) error {
	var countEmail, countUsername, countPhone int

	// Vérifier l'unicité de l'email, mais exclure l'email de l'utilisateur actuel
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", user.Email, userID).Scan(&countEmail)
	if err != nil {
		log.Println("Failed to check email exist:", err)
		return fmt.Errorf("failed to check email exist: %w", err)
	}
	if countEmail > 0 {
		return fmt.Errorf("email already exists")
	}

	// Vérifier l'unicité du nom d'utilisateur, mais exclure celui de l'utilisateur actuel
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", user.Username, userID).Scan(&countUsername)
	if err != nil {
		log.Println("Failed to check username exist:", err)
		return fmt.Errorf("failed to check username exist: %w", err)
	}
	if countUsername > 0 {
		return fmt.Errorf("username already exists")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE phone = ? AND id != ?", user.Phone, userID).Scan(&countPhone)
	if err != nil {
		log.Println("Failed to check phone exist:", err)
		return fmt.Errorf("failed to check phone exist: %w", err)
	}
	if countPhone > 0 {
		log.Println("Phone already exists")
		return fmt.Errorf("phone already exists")
	}
	return nil
}

func IsValidPhoneNumber(phone string) bool {
	var validPhoneNumber = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return validPhoneNumber.MatchString(phone)
}
