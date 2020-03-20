package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	db     *sql.DB
	router *http.ServeMux
}

func main() {
	db, err := sql.Open("mysql", os.Getenv("MYSQL_DSN"))

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	router := http.NewServeMux()

	s := server{db, router}

	s.routes()

	log.Fatal(http.ListenAndServe(":8080", s.router))
}
