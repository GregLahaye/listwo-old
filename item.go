package main

import (
	"encoding/json"
	"net/http"
)

type item struct {
	ID       string `json:"id"`
	Position string `json:"position"`
	Title    string `json:"title"`
}

func (s *server) ownsColumn(userID, columnID string) bool {
	row := s.db.QueryRow("SELECT `uuid` FROM `User` WHERE `id` = (SELECT `user_id` FROM `List` WHERE `id` = (SELECT `list_id` FROM `Column` WHERE uuid = ?))", columnID)

	var ownerID string

	err := row.Scan(&ownerID)

	if err != nil {
		return false
	}

	return userID == ownerID
}

func (s *server) handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetItems(w, r)
	case http.MethodPost:
		s.handleCreateItem(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetItems(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	values := r.URL.Query()

	columnID := values.Get("column")

	if columnID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !s.ownsColumn(userID, columnID) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	rows, err := s.db.Query("SELECT `uuid`, `position`, `title` FROM `Item` WHERE `column_id` = (SELECT `id` FROM `Column` WHERE `uuid` = ?)", columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []item

	for rows.Next() {
		var item item

		err := rows.Scan(&item.ID, &item.Position, &item.Title)

		if err != nil {
			break
		}

		items = append(items, item)
	}

	err = rows.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func (s *server) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	columnID := r.FormValue("column")
	title := r.FormValue("title")

	if columnID == "" || title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	if !s.ownsColumn(userID, columnID) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO `Item` (`uuid`, `position`, `title`, `column_id`) VALUES (UUID(), 0, ?, (SELECT `id` FROM `Column` WHERE `uuid` = ?))")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := stmt.Exec(title, columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err = s.db.Prepare("UPDATE `Item` SET `position` = `position` + 1 WHERE `column_id` = (SELECT `id` FROM `Column` WHERE `uuid` = ?)")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemPK, err := result.LastInsertId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := s.db.QueryRow("SELECT `uuid` FROM `Item` WHERE `id` = ?", itemPK)

	var itemID string

	err = row.Scan(&itemID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itemID)
}
