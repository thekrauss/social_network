package controllers

func (s *MyServer) routes() {
	//	s.Router.Handle("/register", Chain(s.RegisterHandler(), LogRequestMiddleware))
	s.Router.Handle("/login", Chain(s.LoginHandler(), LogRequestMiddleware))
	s.Router.Handle("/create_post", Chain(s.PostHandlers(), LogRequestMiddleware))
	s.Router.Handle("/list_post", Chain(s.ListPostHandler(), LogRequestMiddleware))
	s.Router.Handle("/create_comment", Chain(s.CreateCommentHandler(), LogRequestMiddleware))
	s.Router.Handle("/list_comment", Chain(s.ListCommentHandler(), LogRequestMiddleware))
	s.Router.Handle("/like_comment", Chain(s.LikeCommentHandler(), LogRequestMiddleware))
}
