package controllers

func (s *MyServer) routes() {
	//s.Router.Handle("/frontEnd/", http.StripPrefix("/frontEnd/", http.FileServer(http.Dir("frontEnd"))))
	s.Router.Handle("/register", Chain(s.RegisterHandler(), LogRequestMiddleware))
}
