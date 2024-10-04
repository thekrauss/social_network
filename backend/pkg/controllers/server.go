package controllers

import (
	"backend/pkg/db"
	"context"
	"fmt"
	"log"
	"net/http"
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
	Store  db.Store       // instance de la base de données
	Router *http.ServeMux // routeur HTTP
	Server *http.Server   // serveur HTTP
	//WebSocketChat     *wsk.WebsocketChat // Gestionnaire de chat WebSocket
	GoogleOAuthConfig *oauth2.Config // Configuration OAuth pour Google
	GitHubOAuthConfig *oauth2.Config // Configuration OAuth pour GitHub
}

// Fonction pour créer une nouvelle instance de MyServer
func NewServer(store db.Store /*, wsChat *wsk.WebsocketChat*/) *MyServer {

	router := http.NewServeMux() // Initialisation du routeur HTTP

	// Création de la nouvelle instance de MyServer avec les configurations nécessaires
	server := &MyServer{
		Store:  store,
		Router: router,
		//WebSocketChat: wsChat,
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

	server.routes() // initialisation des routes du serveur

	fmt.Println(ColorBlue, "(http://localhost:8079) - Server started on port", port, ColorReset)
	fmt.Println(ColorGreen, "[SERVER_INFO] : To stop the server : Ctrl + c", ColorReset)

	// Configuration du serveur HTTP
	srv := &http.Server{
		Addr:              "localhost:8080", // adresse à laquelle le serveur écoute
		Handler:           router,           // routeur pour gérer les requêtes
		ReadHeaderTimeout: 15 * time.Second, // délai d'attente pour lire l'en-tête
		ReadTimeout:       15 * time.Second, // délai d'attente pour lire le corps de la requête
		WriteTimeout:      10 * time.Second, // délai d'attente pour écrire la réponse
		IdleTimeout:       30 * time.Second, // délai d'attente pour les connexions inactives
	}

	server.Server = srv // assignation de l'instance du serveur à MyServer

	return server
}

// fonction pour arrêter le serveur proprement avec gestion du contexte
func (s *MyServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// middleware pour logger les requêtes HTTP et appel du prochain handler
func LogRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v], %v", r.Method, r.RequestURI)
		next(w, r)
	}
}

// fonction Chain pour empiler les middlewares
func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		f = middleware(f)
	}
	return f
}
