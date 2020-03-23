package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func testHandler(method, target string, form url.Values) *httptest.ResponseRecorder {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_MOCK_DSN"))

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	router := http.NewServeMux()

	s := server{db, router}

	s.routes()

	r := httptest.NewRequest(method, target, strings.NewReader(form.Encode()))

	r.Form = form

	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, r)

	return w
}

func assert(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
