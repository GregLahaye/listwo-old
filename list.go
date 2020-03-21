package main

import (
	"encoding/json"
	"net/http"
)

type list struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *server) handleList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetList(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetList(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	values := r.URL.Query()

	listID := values.Get("list")

	if listID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !s.ownsList(userID, listID) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	row := s.db.QueryRow("SELECT `uuid`, `title` FROM `List` WHERE `uuid` = ?", listID)

	var list list

	err = row.Scan(&list.ID, &list.Title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func (s *server) handleLists(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetLists(w, r)
	case http.MethodPost:
		s.handleCreateList(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetLists(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	rows, err := s.db.Query("SELECT `uuid`, `title` FROM `List` WHERE `user_id` = (SELECT `id` FROM `User` WHERE `uuid` = ?)", userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var lists []list

	for rows.Next() {
		var list list

		err := rows.Scan(&list.ID, &list.Title)

		if err != nil {
			break
		}

		lists = append(lists, list)
	}

	err = rows.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lists)
}

func (s *server) handleCreateList(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")

	if title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO `List` (`uuid`, `title`, `user_id`) VALUES (UUID(), ?, (SELECT `id` FROM `User` WHERE `uuid` = ?))")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	result, err := stmt.Exec(title, userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listPK, err := result.LastInsertId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := s.db.QueryRow("SELECT `uuid` FROM `List` WHERE `id` = ?", listPK)

	var listID string

	err = row.Scan(&listID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listID)
}
