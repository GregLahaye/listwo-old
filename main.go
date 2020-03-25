package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type server struct {
	db     *sql.DB
	router *http.ServeMux
}

var (
	dbUser                 = os.Getenv("DB_USER")
	dbPwd                  = os.Getenv("DB_PASS")
	instanceConnectionName = os.Getenv("INSTANCE_CONNECTION_NAME")
	dbName                 = os.Getenv("DB_NAME")
)

func main() {
	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s", dbUser, dbPwd, instanceConnectionName, dbName)

	db, err := sql.Open("mysql", dbURI)

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
