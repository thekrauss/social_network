package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s MyServer) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		// Limiter la taille du formulaire multipart (20 MB max ici)
		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form RegisterHandler", http.StatusBadRequest)
			return
		}

		// upload de l'image (avatar) si fourni
		avatarPath, err := UploadImages(w, r, "./uploads/avatars")
		if err != nil && err != http.ErrMissingFile {
			log.Println("Failed to upload avatar:", err)
			http.Error(w, "Failed to upload avatar", http.StatusBadRequest)
			return
		}

		log.Println("Form data received RegisterHandler:", r.FormValue("data"))

		var user models.User
		if err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&user); err != nil {
			log.Println("Failed to decode request payload:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user.Username = strings.TrimSpace(user.Username)
		if len(user.Username) < 3 || len(user.Username) > 30 {
			http.Error(w, "Username must be between 3 and 30 characters", http.StatusBadRequest)
			return
		}

		user.FirstName = strings.TrimSpace(user.FirstName)
		user.LastName = strings.TrimSpace(user.LastName)
		if len(user.FirstName) < 2 || len(user.FirstName) > 30 || len(user.LastName) < 2 || len(user.LastName) > 30 {
			http.Error(w, "FirstName and LastName must be between 2 and 30 characters", http.StatusBadRequest)
			return
		}

		if !IsValidGender(user.Gender) {
			http.Error(w, "Gender must be 'Homme' or 'Femme'.", http.StatusBadRequest)
			return
		}

		if len(user.Password) < 6 || len(user.Password) > 16 {
			http.Error(w, "Password must be between 6 and 16 characters long", http.StatusBadRequest)
			return
		}

		if avatarPath != "" {
			user.Avatar = avatarPath
		} else {
			user.Avatar = ""
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			log.Println("Failed to open database:", err)
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

		log.Println("Database and table ready for register")

		// enregistrement de l'utilisateur
		if err := RegisterUser(w, r, DB, user); err != nil {
			log.Println("Failed to create user:", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		response := models.Response{
			Message: "User registered successfully",
			User:    user,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println("Failed to encode response:", err)
			http.Error(w, "Failed to send response", http.StatusInternalServerError)
			return
		}
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request, DB *sql.DB, user models.User) error {

	if !IsValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Hashage du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Insertion de l'utilisateur
	err = CreateUser(DB, user)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(db *sql.DB, user models.User) error {
	var countEmail, countUsername int

	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", user.Email).Scan(&countEmail)
	if err != nil {
		log.Println("Failed to check email existence:", err)
		return fmt.Errorf("failed to check email existence: %w", err)
	}
	if countEmail > 0 {
		log.Println("Email already exists")
		return fmt.Errorf("email already exists")
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&countUsername)
	if err != nil {
		log.Println("Failed to check username existence:", err)
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if countUsername > 0 {
		log.Println("Username already exists")
		return fmt.Errorf("username already exists")
	}

	DateOfBirth, err := time.Parse("2006-01-02", user.DateOfBirth)
	if err != nil {
		log.Println("Invalid date format:", err)
		return fmt.Errorf("invalid date format: %w", err)
	}

	// Générer un UUID pour l'utilisateur
	userID := uuid.Must(uuid.NewV4())

	query := `INSERT INTO users 
	(id, username, age, email, password_hash, first_name, last_name, role, gender, date_of_birth, avatar, bio, phone_number, address, is_private, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Exec(query, userID, user.Username, user.Age, user.Email, user.Password, user.FirstName, user.LastName, user.Role, user.Gender, DateOfBirth, sql.NullString{String: user.Avatar, Valid: user.Avatar != ""}, sql.NullString{String: user.Bio, Valid: user.Bio != ""}, sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""}, sql.NullString{String: user.Address, Valid: user.Address != ""}, user.IsPrivate, time.Now(), time.Now())

	if err != nil {
		log.Println("Failed to execute insert query:", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	log.Println("User successfully created")
	return nil
}

func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}

func IsValidGender(gender string) bool {
	return gender == "Homme" || gender == "Femme"
}
