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

	port := os.Getenv("PORT")
	addr := ":" + port

	log.Fatal(http.ListenAndServe(addr, s.router))
}
