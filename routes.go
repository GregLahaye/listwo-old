package main

import (
	"net/http"
	"os"
)

func (s *server) routes() {
	s.router.HandleFunc("/signup", middleware(s.handleSignUp))
	s.router.HandleFunc("/signin", middleware(s.handleSignIn))
	s.router.HandleFunc("/lists", middleware(s.handleLists))
	s.router.HandleFunc("/list", middleware(s.handleList))
	s.router.HandleFunc("/columns", middleware(s.handleColumns))
	s.router.HandleFunc("/items", middleware(s.handleItems))
}

func middleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", os.Getenv("LISTWO_API_ALLOW_ORIGIN"))
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		h(w, r)
	}
}
