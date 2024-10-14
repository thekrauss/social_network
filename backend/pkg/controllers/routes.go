package controllers

func (s *MyServer) routes() {

	s.Router.Handle("/register", Chain(s.RegisterHandler(), LogRequestMiddleware))
	s.Router.Handle("/login", Chain(s.LoginHandler(), LogRequestMiddleware))
	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/create_post", Chain(s.PostHandlers(), LogRequestMiddleware))
	s.Router.Handle("/list_post", Chain(s.ListPostHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/create_comment", Chain(s.CreateCommentHandler(), LogRequestMiddleware))
	s.Router.Handle("/list_comment", Chain(s.ListCommentHandler(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/like_post", Chain(s.LikePost(), LogRequestMiddleware))
	s.Router.Handle("/unlike_post", Chain(s.UnlikePost(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/like_comment", Chain(s.LikeComment(), LogRequestMiddleware))
	s.Router.Handle("/unlike_comment", Chain(s.UnLikeComment(), LogRequestMiddleware))

	/*-------------------------------------------------------------------------------*/

	s.Router.Handle("/update_profil", Chain(s.UpdateProfileHandler(), LogRequestMiddleware))

}
