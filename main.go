package main

import (
	"fmt"
	"log"
	"net/http"
	"usuarios-crud/database"
	"usuarios-crud/handlers"
)

func main() {
	// Abre a conexão com o SQLite e cria a tabela na primeira execução.
	err := database.InitDB("./users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Cada rota principal delega para o handler correto conforme o método HTTP.
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/user", handleUser)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// /users concentra criação e listagem.
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.CreateUserHandler(w, r)
	case http.MethodGet:
		handlers.ListUsersHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// /user concentra busca por id, atualização e remoção.
func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handlers.GetUserHandler(w, r)
	case http.MethodPut:
		handlers.UpdateUserHandler(w, r)
	case http.MethodDelete:
		handlers.DeleteUserHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
