package api

import (
    "fmt"
    "log"
    "encoding/json"
    "database/sql"
    "net/http"
    _ "github.com/lib/pq"
)

type Quote struct {
    Id     int
    Quote  string
    Author string
}

type Server struct {
    Db *sql.DB
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
    query := "SELECT quote, author FROM quote"
    rows, err := s.Db.Query(query)
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

func (s *Server) HandlePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var quote Quote
        decoder := json.NewDecoder(r.Body)
        decodeError := decoder.Decode(&quote)
        if decodeError != nil {
            log.Println("[ERROR]", decodeError)
        }
        query := fmt.Sprintf("INSERT INTO quote (id, quote, author) VALUES (%d, '%s', '%s')", quote.Id, quote.Quote, quote.Author)
        _, queryRrror := s.Db.Exec(query)
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
