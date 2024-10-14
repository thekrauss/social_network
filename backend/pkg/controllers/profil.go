package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gofrs/uuid"
)

func (s *MyServer) UpdateProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
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

			if updatedUser.FirstName != "" && (len(updatedUser.FirstName) < 4 || len(updatedUser.FirstName) > 16) {
				http.Error(w, "First name must be between 4 and 16 characters", http.StatusBadRequest)
				return
			}

			if updatedUser.LastName != "" && (len(updatedUser.LastName) < 4 || len(updatedUser.LastName) > 16) {
				http.Error(w, "Last name must be between 4 and 16 characters", http.StatusBadRequest)
				return
			}

			if updatedUser.Gender != "" && !IsValidGender(updatedUser.Gender) {
				http.Error(w, "Invalid gender. Must be 'Homme', 'Femme', or 'other'", http.StatusBadRequest)
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

			if updatedUser.DateOfBirth != "" {
				http.Error(w, "Date of birth cannot be modified", http.StatusBadRequest)
				return
			}

			err = checkUser(DB, updatedUser)
			if err != nil {
				log.Println("Failed to check user:", err)
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

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

		} else {
			http.NotFound(w, r)
		}
	}
}

func IsValidPhoneNumber(phone string) bool {
	var validPhoneNumber = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return validPhoneNumber.MatchString(phone)
}

func checkUser(db *sql.DB, user models.User) error {
	var countusername, countEmail, countFirstname, countLastname, countPhone int

	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&countEmail)
	if err != nil {
		log.Println("Failed to check email existence:", err)
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if countEmail > 0 {
		log.Println("Email already exists")
		return fmt.Errorf("email already exists")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE first_name = ?", user.Username).Scan(&countusername)
	if err != nil {
		log.Println("Failed to check first name existence:", err)
		return fmt.Errorf("failed to check first name existence: %w", err)
	}
	if countusername > 0 {
		log.Println("Failed to check username existence:", err)
		return fmt.Errorf("failed to check username existence: %w", err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE first_name = ?", user.FirstName).Scan(&countFirstname)
	if err != nil {
		log.Println("Failed to check first name existence:", err)
		return fmt.Errorf("failed to check first name existence: %w", err)
	}
	if countFirstname > 0 {
		log.Println("First name already exists")
		return fmt.Errorf("first name already exists")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE last_name = ?", user.LastName).Scan(&countLastname)
	if err != nil {
		log.Println("Failed to check last name existence:", err)
		return fmt.Errorf("failed to check last name existence: %w", err)
	}
	if countLastname > 0 {
		log.Println("Last name already exists")
		return fmt.Errorf("last name already exists")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE phone = ?", user.Phone).Scan(&countPhone)
	if err != nil {
		log.Println("Failed to check phone existence:", err)
		return fmt.Errorf("failed to check phone existence: %w", err)
	}
	if countPhone > 0 {
		log.Println("Phone already exists")
		return fmt.Errorf("phone already exists")
	}

	return nil
}
