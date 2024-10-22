package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Exp      int64     `json:"exp"`
}

var jwtKey = []byte("my_secret_key")

type contextKey string

const userIDKey contextKey = "userID"

// GenerateJWT génère un JWT avec UserID, Username et une date d'expiration
func GenerateJWT(userID uuid.UUID, username string) (string, error) {
	// Générer le header JWT
	header := strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`)), "=")

	// Générer les claims avec expiration
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Exp:      expirationTime,
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payload := strings.TrimRight(base64.URLEncoding.EncodeToString(claimsJSON), "=")

	// Générer la signature HMAC
	signature := generateHMACSHA256(header + "." + payload)
	token := fmt.Sprintf("%s.%s.%s", header, payload, signature)

	return token, nil
}

// generateHMACSHA256 génère la signature HMAC-SHA256
func generateHMACSHA256(data string) string {
	h := hmac.New(sha256.New, jwtKey)
	h.Write([]byte(data))
	return strings.TrimRight(base64.URLEncoding.EncodeToString(h.Sum(nil)), "=")
}

// VerifyJWT vérifie le token JWT et retourne les claims si valides
func VerifyJWT(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		log.Println("Invalid token format")
		return nil, fmt.Errorf("invalid token format")
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	expectedSignature := generateHMACSHA256(header + "." + payload)
	if signature != expectedSignature {
		log.Println("Invalid token signature")
		return nil, fmt.Errorf("invalid token signature")
	}

	payloadData, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		log.Println("Invalid token payload")
		return nil, fmt.Errorf("invalid token payload")
	}

	var claims Claims
	err = json.Unmarshal(payloadData, &claims)
	if err != nil {
		log.Println("Invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	// Vérifie si le token a expiré
	if time.Now().Unix() > claims.Exp {
		log.Println("Token has expired")
		return nil, fmt.Errorf("token has expired")
	}

	log.Println("Token is valid, UserID:", claims.UserID)
	return &claims, nil
}

func (s *MyServer) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("No Authorization header")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Println("Invalid Authorization format")
			http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			return
		}

		token := tokenParts[1]
		claims, err := VerifyJWT(token)
		if err != nil {
			log.Println("Token verification failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Injecter l'ID utilisateur dans le contexte
		log.Println("User ID from token:", claims.UserID)
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
