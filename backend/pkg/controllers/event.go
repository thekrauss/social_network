package controllers

import (
	"backend/pkg/models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreateEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'utilisateur connecté
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		var event models.GroupEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		event.ID = uuid.Must(uuid.NewV4())
		event.UserID = userID

		// Ouvrir la base de données
		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		tx, err := DB.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
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

		if event.Title == "" || event.Description == "" || event.EventDate.IsZero() {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Insérer l'événement dans la base de données
		query := `INSERT INTO group_events (id, group_id, user_id, title, description, event_date) VALUES (?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(query, event.ID, event.GroupID, event.UserID, event.Title, event.Description, event.EventDate)
		if err != nil {
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Event created successfully"))
	}
}

func (s *MyServer) ListEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		query := r.URL.Query()
		GroupIDStr := query.Get("group_id")
		if GroupIDStr == "" {
			http.Error(w, "Group ID not provided", http.StatusBadRequest)
			return
		}

		groupID, err := uuid.FromString(GroupIDStr)
		if err != nil {
			http.Error(w, "Invalid Group ID", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
		defer DB.Close()

		// récupération des paramètres de pagination : page et limit
		page := 1
		limit := 10

		// Parse des paramètres page et limit
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

		offset := (page - 1) * limit
		log.Printf("Fetching event from database (page: %d, limit: %d)\n", page, limit)

		events, err := GetEventByGroup(DB, groupID, offset, limit)
		if err != nil {
			http.Error(w, "Failed to retrDBieve event", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(events); err != nil {
			http.Error(w, "Failed to encode events", http.StatusInternalServerError)
		}

	}
}

func (s *MyServer) InviteToEventHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			log.Println("User not logged in", ok)
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		var response models.EventResponse
		if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		response.UserID = userID

		if response.Response != "Going" && response.Response != "NotGoing" {
			log.Println("Invalid response")
			http.Error(w, "Invalid response", http.StatusBadRequest)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}

		tx, err := DB.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
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

		var query string

		if response.Response == "Going" {
			query = "UPDATE group_events SET options.going = options.going + 1 WHERE id = ?"
		} else {
			query = "UPDATE group_events SET options.not_going = options.not_going + 1 WHERE id = ?"

		}

		_, err = tx.Exec(query, response.EventID)
		if err != nil {
			log.Println("Failed to update event response")
			http.Error(w, "Failed to update event response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Response recorded successfully"))

	}
}

func GetEventByGroup(DB *sql.DB, groupID uuid.UUID, offset, limit int) ([]models.GroupEvent, error) {

	if DB == nil {
		return nil, errors.New("database connection is nil")
	}
	rows, err := DB.Query("SELECT id, group_id, user_id, title, description, event_date, created_at FROM group_events WHERE group_id = ? LIMIT ? OFFSET ?", groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.GroupEvent
	for rows.Next() {
		var event models.GroupEvent
		err := rows.Scan(&event.ID, &event.GroupID, &event.UserID, &event.Title, &event.Description, &event.EventDate, &event.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Par défaut, event.Options est à zéro, ou récupère des données  si nécessaire
		event.Options = models.EventOptions{
			Going:    0,
			NotGoing: 0,
		}

		events = append(events, event)
	}

	return events, nil
}
