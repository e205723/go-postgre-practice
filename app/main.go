package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"yoshisaur/api"
	_ "github.com/lib/pq"
)

const (
    host     = "db"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "postgres"
)

func main() {
    postqreslInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", postqreslInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    server := api.Server{Db: db}

    http.HandleFunc("/get-user", server.HandleGet)
    http.HandleFunc("/post-user", server.HandlePost)
    log.Fatal(http.ListenAndServe(":8080", nil))
    if err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
