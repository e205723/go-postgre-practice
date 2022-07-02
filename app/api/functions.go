package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("oG+kRyzbqMpuvm2AkRHVhMbLvYoiwVMjs7WtBaxrksxL5Ex646JlJA==")

type User struct {
    Name  string
    Password  string
}

type Server struct {
    Db *sql.DB
}

type Claims struct {
    Name string
    jwt.StandardClaims
}

func (s *Server) HandleGet(w http.ResponseWriter, r *http.Request) {
    query := "SELECT id, name FROM users"
    rows, err := s.Db.Query(query)
    defer rows.Close()
    if err != nil {
        log.Println("[ERROR]", err)
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad"))
        return
    } else {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    }
    for rows.Next() {
        var (
            id int
            name string
        )
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal(err)
            return
        }
        log.Printf("id: '%d'\n user: %s\n\n", id, name)
    }
}

func (s *Server) HandlePost(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var user User
        decoder := json.NewDecoder(r.Body)
        decodeError := decoder.Decode(&user)
        if decodeError != nil {
            log.Println("[ERROR]", decodeError)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("bad"))
            return
        }
        query := fmt.Sprintf("INSERT INTO users (name, password) VALUES ('%s', '%s')", user.Name, user.Password)
        _, queryError := s.Db.Exec(query)
        if queryError != nil {
            log.Println("[ERROR]", queryError)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("bad"))
            return
        } else {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte("ok"))
        }
    } else {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad"))
        return
    }
}

func (s *Server) Signin(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var user User
        decoder := json.NewDecoder(r.Body)
        decodeError := decoder.Decode(&user)
        if decodeError != nil {
            log.Println("[ERROR]", decodeError)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("bad"))
            return
        }
        queryToGetPassword := fmt.Sprintf("SELECT password FROM users WHERE name='%s'", user.Name)
        rows, queryError := s.Db.Query(queryToGetPassword)
        defer rows.Close()
        if queryError != nil {
            log.Println("[ERROR]", queryError)
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("bad"))
            return
        }
        var correctPassword string
        for rows.Next() {
            if err := rows.Scan(&correctPassword); err != nil {
                log.Fatal(err)
                return
            }
        }
        correctPassword = strings.TrimRight(correctPassword, " ")
        if user.Password != correctPassword {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        expirationTime := time.Now().Add(24 * time.Hour)
        claims := &Claims{
            Name: user.Name,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        cookie := &http.Cookie{
            Name: "token",
            Value: tokenString,
            Expires: expirationTime,
        }
        http.SetCookie(w, cookie)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    } else {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("bad"))
        return
    }
}

func (s *Server) Welcome(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tknStr := c.Value
    claims := &Claims{}
    tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if !tkn.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(fmt.Sprintf("Welcome %s", claims.Name)))
}

func (s *Server) Refresh(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    tknStr := c.Value
    claims := &Claims{}
    tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if !tkn.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
    if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 24*time.Hour {
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    expirationTime := time.Now().Add(24 * time.Hour)
    claims.ExpiresAt = expirationTime.Unix()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    r.AddCookie(&http.Cookie{
        Name: "token",
        Value: tokenString,
        Expires: expirationTime,
    })
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}
