package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"s-network/pkg/db"
	"s-network/pkg/wsk"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const (
	ColorGreen = "\033[32m"
	ColorBlue  = "\033[34m"
	ColorReset = "\033[0m"
	port       = ":8080"
)

// Structure pour le serveur
type MyServer struct {
	Store             db.Store           // Instance de la base de données
	Router            *http.ServeMux     // Routeur HTTP
	Server            *http.Server       // Serveur HTTP
	WebSocketChat     *wsk.WebsocketChat // Gestionnaire de chat WebSocket
	GoogleOAuthConfig *oauth2.Config     // Configuration OAuth pour Google
	GitHubOAuthConfig *oauth2.Config     // Configuration OAuth pour GitHub
}

// Fonction pour créer une nouvelle instance de MyServer
func NewServer(store db.Store, wsChat *wsk.WebsocketChat) *MyServer {

	router := http.NewServeMux() // Initialisation du routeur HTTP

	// Création de la nouvelle instance de MyServer avec les configurations nécessaires
	server := &MyServer{
		Store:         store,
		Router:        router,
		WebSocketChat: wsChat,
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     "your-google-client-id",
			ClientSecret: "your-google-client-secret",
			RedirectURL:  "http://localhost:8079/auth/google/callback",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		GitHubOAuthConfig: &oauth2.Config{
			ClientID:     "your-github-client-id",
			ClientSecret: "your-github-client-secret",
			RedirectURL:  "http://localhost:8079/auth/github/callback",
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
	}

	server.routes() // Initialisation des routes du serveur

	fmt.Println(ColorBlue, "(http://localhost:8079) - Server started on port", port, ColorReset)
	fmt.Println(ColorGreen, "[SERVER_INFO] : To stop the server : Ctrl + c", ColorReset)

	// Configuration du serveur HTTP
	srv := &http.Server{
		Addr:              "localhost:8080", // Adresse à laquelle le serveur écoute
		Handler:           router,           // Routeur pour gérer les requêtes
		ReadHeaderTimeout: 15 * time.Second, // Délai d'attente pour lire l'en-tête
		ReadTimeout:       15 * time.Second, // Délai d'attente pour lire le corps de la requête
		WriteTimeout:      10 * time.Second, // Délai d'attente pour écrire la réponse
		IdleTimeout:       30 * time.Second, // Délai d'attente pour les connexions inactives
	}

	server.Server = srv // Assignation de l'instance du serveur à MyServer

	return server // Retourne la nouvelle instance de MyServer
}

// Fonction pour arrêter le serveur proprement
func (s *MyServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx) // Arrêt du serveur avec gestion du contexte
}

// Middleware pour logger les requêtes HTTP
func LogRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Journalisation de la méthode et de l'URI de la requête
		log.Printf("[%v], %v", r.Method, r.RequestURI)
		next(w, r) // Appel du prochain handler
	}
}

// Fonction Chain pour empiler les middlewares
func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		f = middleware(f)
	}
	return f // Retourne la fonction handler avec les middlewares appliqués
}
