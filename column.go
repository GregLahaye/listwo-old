package main

import (
	"encoding/json"
	"net/http"
)

type column struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *server) ownsList(userID, listID string) bool {
	row := s.db.QueryRow("SELECT `uuid` FROM `User` WHERE `id` = (SELECT `user_id` FROM `List` WHERE `uuid` = ?)", listID)

	var ownerID string

	err := row.Scan(&ownerID)

	if err != nil {
		return false
	}

	return userID == ownerID
}

func (s *server) handleColumns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetColumns(w, r)
	case http.MethodPost:
		s.handleCreateColumn(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetColumns(w http.ResponseWriter, r *http.Request) {
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

	rows, err := s.db.Query("SELECT `uuid`, `title` FROM `Column` WHERE `list_id` = (SELECT `id` FROM `List` WHERE `uuid` = ?)", listID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var columns []column

	for rows.Next() {
		var column column

		err := rows.Scan(&column.ID, &column.Title)

		if err != nil {
			break
		}

		columns = append(columns, column)
	}

	err = rows.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(columns)
}

func (s *server) handleCreateColumn(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	listID := r.FormValue("list")
	title := r.FormValue("title")

	if listID == "" || title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if !s.ownsList(userID, listID) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO `Column` (`uuid`, `title`, `list_id`) VALUES (UUID(), ?, (SELECT `id` FROM `List` WHERE `uuid` = ?))")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := stmt.Exec(title, listID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	columnPK, err := result.LastInsertId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := s.db.QueryRow("SELECT `uuid` FROM `Column` WHERE `id` = ?", columnPK)

	var columnID string

	err = row.Scan(&columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(columnID)
}