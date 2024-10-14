package controllers

import (
	"net/http"
)

func (s *MyServer) OnlineUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *MyServer) MessagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
