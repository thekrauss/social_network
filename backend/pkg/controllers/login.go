package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponses struct {
	Token   string `json:"token,omitempty"`
	Message string `json:"message"`
}

func (s MyServer) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Println("Failed to parse form", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			DB, err := s.Store.OpenDatabase()
			if err != nil {
				log.Println("Failed to open database", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			defer DB.Close()

			fmt.Println("ping to database successful")

			identifier := r.FormValue("identifier")
			password := r.FormValue("password")

			log.Println("Login attempt with identifier:", identifier)
			log.Println("Password provided:", password)

			if identifier == "" || password == "" {
				log.Println("Identifier or password is empty")
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}

			if len(identifier) < 6 || len(identifier) > 16 || len(password) < 6 || len(password) > 16 {
				log.Println("Identifier and Password must be between 6 and 16 characters long")
				http.Error(w, "Identifier and Password must be between 6 and 16 characters long", http.StatusBadRequest)
				return
			}

			userID := uuid.Must(uuid.NewV4())
			var storedPassword, username string

			if strings.Contains(identifier, "@") {
				log.Println("Identified as email")
				userID, err = GetUserIDbyEmail(DB, identifier)
				if err != nil {
					log.Println("Email does not exist", err)
					http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
					return
				}
				storedPassword, err = GetPasswordByEmail(DB, identifier)
				if err != nil {
					log.Println("Failed to retrieve password by email", err)
					http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
					return
				}

				username, err = GetUsernameByEmail(DB, identifier)
				if err != nil {
					log.Println("Failed to retrieve username by email", err)
					http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
					return
				}

			} else {
				log.Println("Identified as username")
				userID, err = GetUserIDbyUsername(DB, identifier)
				if err != nil {
					if err == sql.ErrNoRows {
						log.Println("Username does not exist:", identifier)
					} else {
						log.Println("Failed to retrieve user ID by username:", err)
					}
					http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
					return
				}
				storedPassword, err = GetPasswordByUsername(DB, identifier)
				if err != nil {
					log.Println("Failed to retrieve password by username", err)
					http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
					return
				}
				username = identifier
			}

			log.Println("UserID:", userID, "StoredPassword:", storedPassword)

			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
			if err != nil {
				log.Println("Incorrect password", err)
				http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
				return
			}

			token, err := GenerateJWT(userID, username)
			if err != nil {
				log.Println("Failed to generate token", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   token,
				Expires: time.Now().Add(24 * time.Hour),
			})

			// Set username cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "username",
				Value:   username,
				Expires: time.Now().Add(24 * time.Hour),
			})

			log.Println("User logged in successfully, userID:", userID)
			SendJSONResponse(w, LoginResponses{Token: token, Message: "Login successful"}, http.StatusOK)

		} else {
			http.NotFound(w, r)
		}
	}
}

func (s *MyServer) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Now(),
		})
		w.Write([]byte("Logged out successfully"))
	}
}

func SendJSONResponse(w http.ResponseWriter, response LoginResponses, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
