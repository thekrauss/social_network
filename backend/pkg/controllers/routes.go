package controllers

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func (s *MyServer) routes() {

	s.Router.HandleFunc("/protected", Chain(s.ProtectedHandler(), LogRequestMiddleware, s.Authenticate))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/register", Chain(s.RegisterHandler(), LogRequestMiddleware))
	s.Router.Handle("/login", Chain(s.LoginHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/logout", Chain(s.LogoutHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.HandleFunc("/auth/github/login-form", Chain(s.GitHubLoginHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/callback", Chain(s.GitHubCallbackHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/register-form", Chain(s.GitHubRegisterHandler(), LogRequestMiddleware))
	s.Router.HandleFunc("/auth/github/callback/register", Chain(s.GitHubCallbackRegisterHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/create_post", Chain(s.CreatePostHandlers(), LogRequestMiddleware, s.Authenticate))
	s.Router.Handle("/list_post", Chain(s.ListPostHandler(), LogRequestMiddleware, s.Authenticate))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/create_comment", Chain(s.CreateCommentHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.Handle("/list_comment", Chain(s.ListCommentHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/like_post", Chain(s.LikePost(), LogRequestMiddleware))
	s.Router.Handle("/unlike_post", Chain(s.UnlikePost(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/like_comment", Chain(s.LikeComment(), LogRequestMiddleware))
	s.Router.Handle("/unlike_comment", Chain(s.UnLikeComment(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/view_profil", Chain(s.GetUserProfil(), LogRequestMiddleware))
	s.Router.Handle("/my_profil", Chain(s.MyProfil(), LogRequestMiddleware))
	s.Router.Handle("/update_profil", Chain(s.UpdateProfileHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/list_group", Chain(s.ListGroupsHandler(), LogRequestMiddleware))
	s.Router.Handle("/create_group", Chain(s.CreateGroupHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.Handle("/invit_group", Chain(s.InviteToGroupHandler(), LogRequestMiddleware))
	s.Router.Handle("/create_post_group", Chain(s.CreatePostGroupHandler(), LogRequestMiddleware, s.Authenticate))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/create_event", Chain(s.CreateEventHandler(), LogRequestMiddleware, s.Authenticate))
	s.Router.Handle("/list_event", Chain(s.ListEvent(), LogRequestMiddleware))
	s.Router.Handle("/invit_event", Chain(s.InviteToEventHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

}

func (s *MyServer) ProtectedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(uuid.UUID)
		w.Write([]byte(fmt.Sprintf("Hello, user %d", userID)))
	}
}
