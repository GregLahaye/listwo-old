package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func generateAccessToken() (string, error) {
	return generateRandomString(256)
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)

	return base64.URLEncoding.EncodeToString(b), err
}

func (s *server) getUser(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")

	row := s.db.QueryRow("SELECT uuid FROM User WHERE id = (SELECT user_id FROM Session WHERE access_token = ?)", authorization)

	var userID string

	err := row.Scan(&userID)

	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (s *server) handleSignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(email) < 3 || len(password) < 8 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	hash, err := hashPassword(password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO User (uuid, email, password) VALUES (UUID(), ?, ?)")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(email, hash)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := s.db.QueryRow("SELECT uuid FROM User WHERE email = ?", email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userID string

	err = row.Scan(&userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.handleSignIn(w, r)
}

func (s *server) handleSignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	row := s.db.QueryRow("SELECT uuid, password FROM User WHERE email = ?", email)

	var userID string
	var hash string

	err := row.Scan(&userID, &hash)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	match := checkPasswordHash(password, hash)

	if match != true {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	accessToken, err := generateAccessToken()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO Session (access_token, user_id) VALUES (?, (SELECT id FROM User WHERE uuid = ?))")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(accessToken, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accessToken)
}