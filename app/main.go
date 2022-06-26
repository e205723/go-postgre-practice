package main

import (
    "fmt"
    "log"
    "encoding/json"
    "database/sql"
    "net/http"
    _ "github.com/lib/pq"
)

const (
    host     = "db"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "postgres"
)

type Quote struct {
    Id     int
    Quote  string
    Author string
}

type server struct {
    db *sql.DB
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
    query := "SELECT quote, author FROM quote"
    rows, err := s.db.Query(query)
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad"))
    } else {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    }
    defer rows.Close()
    for rows.Next() {
        var (
            quote string
            author string
        )
        if err := rows.Scan(&quote, &author); err != nil {
            log.Fatal(err)
        }
        log.Printf("Quote: '%s'\n Author: %s\n\n", quote, author)
    }
}

func (s *server) handlePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var quote Quote
        decoder := json.NewDecoder(r.Body)
        decodeError := decoder.Decode(&quote)
        if decodeError != nil {
            log.Println("[ERROR]", decodeError)
        }
        query := fmt.Sprintf("INSERT INTO quote (id, quote, author) VALUES (%d, '%s', '%s')", quote.Id, quote.Quote, quote.Author)
        _, queryRrror := s.db.Exec(query)
        if queryRrror != nil {
            log.Println("[ERROR]", queryRrror)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("bad"))
        } else {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("ok"))
        }
    } else {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad"))
    }
}

func main() {
    postqreslInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", postqreslInfo)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    server := server{db: db}

    http.HandleFunc("/get-quote", server.handleGet)
    http.HandleFunc("/post-quote", server.handlePost)
    log.Fatal(http.ListenAndServe(":8080", nil))
    if err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
