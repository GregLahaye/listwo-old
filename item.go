package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type item struct {
	ID       string `json:"id"`
	Position int    `json:"position"`
	Title    string `json:"title"`
}

func (s *server) ownsItem(userID, itemID string) bool {
	row := s.db.QueryRow("SELECT `uuid` FROM `User` WHERE `id` = (SELECT `user_id` FROM `List` WHERE `id` = (SELECT `list_id` FROM `Column` WHERE id = (SELECT `column_id` FROM `Item` WHERE `uuid` = ?)))", itemID)

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
	case http.MethodDelete:
		s.handleDeleteItem(w, r)
	case http.MethodPatch:
		s.handleUpdateItem(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleGetItems(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

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
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
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
	userID, err := s.getCurrentUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	columnID := r.FormValue("column")
	title := r.FormValue("title")

	if columnID == "" || title == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !s.ownsColumn(userID, columnID) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	stmt, err := s.db.Prepare("INSERT INTO `Item` (`uuid`, `position`, `title`, `column_id`) SELECT UUID(), (SELECT IFNULL(MAX(`position`) + 1, 0) FROM `Item` WHERE `column_id` = (SELECT `id` FROM `Column` WHERE `uuid` = ?)), ?, (SELECT `id` FROM `Column` WHERE `uuid` = ?)")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := stmt.Exec(columnID, title, columnID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	itemPK, err := result.LastInsertId()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	row := s.db.QueryRow("SELECT `uuid`, `position` FROM `Item` WHERE `id` = ?", itemPK)

	var itemID string
	var position int

	err = row.Scan(&itemID, &position)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := item{
		ID:       itemID,
		Position: position,
		Title:    title,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *server) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	itemID := r.FormValue("id")
	dstStr := r.FormValue("dst")

	if itemID == "" || dstStr == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !s.ownsItem(userID, itemID) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	row := s.db.QueryRow("SELECT `position`, `column_id` FROM `Item` WHERE `uuid` = ?", itemID)

	var src int64
	var columnPK string

	err = row.Scan(&src, &columnPK)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := strconv.ParseInt(dstStr, 10, 64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dir := -1
	lower := src
	upper := dst

	if src > dst {
		dir = 1
		lower = dst
		upper = src
	}

	stmt, err := s.db.Prepare("UPDATE `Item` SET `position` = `position` + ? WHERE `column_id` = ? AND `position` >= ? AND `position` <= ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(dir, columnPK, lower, upper)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err = s.db.Prepare("UPDATE `Item` SET `position` = ? WHERE `uuid` = ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(dst, itemID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, err := s.getCurrentUser(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	itemID := r.FormValue("id")

	if itemID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !s.ownsItem(userID, itemID) {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	row := s.db.QueryRow("SELECT `position`, `column_id` FROM `Item` WHERE `uuid` = ?", itemID)

	var position, columnPK string

	err = row.Scan(&position, &columnPK)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := s.db.Prepare("UPDATE `Item` SET `position` = `position` - 1 WHERE `column_id` = ? AND `position` >= ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = stmt.Exec(columnPK, position)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err = s.db.Prepare("DELETE FROM `Item` WHERE `uuid` = ?")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(itemID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(itemID)
}
