package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dirleiflsilva/crud-de-usuarios-adicionando-persistencia/database"
)

func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	mux := http.NewServeMux()

	// GET /users — list users, optional ?search= query param
	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("search")
		users := db.List(search)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("GET /users: encode error: %v", err)
		}
	})

	// POST /users — create a new user
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
			return
		}
		if body.Name == "" || body.Email == "" {
			http.Error(w, `{"error":"name and email are required"}`, http.StatusBadRequest)
			return
		}

		id, err := newUUID()
		if err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		user := database.User{
			ID:        id,
			Name:      body.Name,
			Email:     body.Email,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(user); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			log.Printf("POST /users: encode error: %v", err)
		}
	})

	// PUT /users/{id} — update an existing user
	mux.HandleFunc("PUT /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var body struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
			return
		}

		found, err := db.Update(id, body.Name, body.Email)
		if err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	// DELETE /users/{id} — remove a user
	mux.HandleFunc("DELETE /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		found, err := db.Delete(id)
		if err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("Server running on http://localhost:3333")
	log.Fatal(http.ListenAndServe(":3333", mux))
}
