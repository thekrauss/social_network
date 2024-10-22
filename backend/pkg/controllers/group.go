package controllers

import (
	"backend/pkg/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
)

func (s *MyServer) CreateGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// récupére l'utilisateur connecté
		userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		var group models.Group
		if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		group.ID = uuid.Must(uuid.NewV4())
		group.CreatorID = userID

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

		// insére le groupe dans la base de données
		query := `INSERT INTO groups (id, name, description, creator_id) VALUES (?, ?, ?, ?)`
		_, err = tx.Exec(query, group.ID, group.Name, group.Description, group.CreatorID)
		if err != nil {
			http.Error(w, "Failed to create group", http.StatusInternalServerError)
			return
		}

		// ajouter le créateur comme membre du groupe avec le rôle de "creator"
		query = `INSERT INTO group_members (id, group_id, user_id, status, role) VALUES (?, ?, ?, 'accepted', 'creator')`
		_, err = tx.Exec(query, uuid.Must(uuid.NewV4()), group.ID, group.CreatorID)
		if err != nil {
			http.Error(w, "Failed to add group creator as member", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Group created successfully"))
	}
}

func (s *MyServer) InviteToGroupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Récupérer l'utilisateur connecté
		inviterID, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			http.Error(w, "User not logged in", http.StatusUnauthorized)
			return
		}

		var inviteData struct {
			GroupID   uuid.UUID `json:"group_id"`
			InviteeID uuid.UUID `json:"invitee_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&inviteData); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
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

		// verifie que l'inviteur est bien membre du groupe
		var inviterRole string
		query := `SELECT role FROM group_members WHERE group_id = ? AND user_id = ? AND status = 'accepted'`
		err = tx.QueryRow(query, inviteData.GroupID, inviterID).Scan(&inviterRole)
		if err != nil || inviterRole == "" {
			http.Error(w, "User not authorized to invite to group", http.StatusUnauthorized)
			return
		}

		// ajoute une invitation avec un statut "pending"
		query = `SELECT status FROM group_members WHERE group_id = ? AND user_id = ?`
		var status string
		err = tx.QueryRow(query, inviteData.GroupID, inviteData.InviteeID).Scan(&status)
		if err == nil && status == "pending" {
			http.Error(w, "User already invited to the group", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User invited to group successfully"))
	}
}

func (s *MyServer) ListGroupsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		DB, err := s.Store.OpenDatabase()
		if err != nil {
			http.Error(w, "Failed to open database", http.StatusInternalServerError)
			return
		}
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

		log.Println("Database and table ready")
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
		// Fetch les posts depuis la base de données avec pagination
		log.Printf("Fetching posts from database (page: %d, limit: %d)\n", page, limit)

		// Requête pour récupérer tous les groupes
		query := `SELECT id, name, description FROM groups  LIMIT ? OFFSET ?`
		rows, err := DB.Query(query, limit, offset)
		if err != nil {
			http.Error(w, "Failed to retrieve groups", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var groups []models.Group
		for rows.Next() {
			var group models.Group
			if err := rows.Scan(&group.ID, &group.Name, &group.Description); err != nil {
				http.Error(w, "Failed to scan group", http.StatusInternalServerError)
				return
			}
			groups = append(groups, group)
		}

		// Répondre avec la liste des groupes
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			http.Error(w, "Failed to encode groups", http.StatusInternalServerError)
		}
	}
}

/*-------------------------------------------------------------------------------*/
