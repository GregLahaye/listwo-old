package main

import (
	"net/http"
	"os"
)

func (s *server) routes() {
	s.router.HandleFunc("/signup", handleCORS(s.handleSignUp))
	s.router.HandleFunc("/signin", handleCORS(s.handleSignIn))
	s.router.HandleFunc("/lists", handleCORS(s.handleLists))
	s.router.HandleFunc("/list", handleCORS(s.handleList))
	s.router.HandleFunc("/columns", handleCORS(s.handleColumns))
	s.router.HandleFunc("/items", handleCORS(s.handleItems))
}

func handleCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if os.Getenv("LISTWO_ENV") == "development" {
			w.Header().Add("Access-Control-Allow-Origin", "*")
		}

		h(w, r)
	}
}
