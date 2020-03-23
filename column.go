package main

import (
	"encoding/json"
	"net/http"
)

type column struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func (s *server) ownsColumn(userID, columnID string) bool {
	row := s.db.QueryRow("SELECT `uuid` FROM `User` WHERE `id` = (SELECT `user_id` FROM `List` WHERE `id` = (SELECT `list_id` FROM `Column` WHERE `uuid` = ?))", columnID)

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
	case http.MethodDelete:
		s.handleDeleteColumn(w, r)
	case http.MethodPatch:
		s.handleUpdateColumn(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetColumns(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

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
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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
	userID, err := s.getCurrentUser(r)

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
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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

	response := column{
		ID:    columnID,
		Title: title,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *server) handleDeleteColumn(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	columnID := r.FormValue("id")

	if !s.ownsColumn(userID, columnID) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	stmt, err := s.db.Prepare("DELETE FROM `Column` WHERE `uuid` = ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(columnID)
}

func (s *server) handleUpdateColumn(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	columnID := r.FormValue("id")
	title := r.FormValue("title")

	if !s.ownsColumn(userID, columnID) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	stmt, err := s.db.Prepare("UPDATE `Column` SET `title` = ? WHERE `uuid` = ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(title, columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(columnID)
}
