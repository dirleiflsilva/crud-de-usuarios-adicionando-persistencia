package database

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"
)

const dbPath = "database/db.json"

// User represents a user entity.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Database holds the in-memory state and manages file persistence.
type Database struct {
	mu    sync.RWMutex
	users []User
}

// New loads the database from disk, or creates an empty one if the file doesn't exist.
func New() (*Database, error) {
	db := &Database{}
	if err := db.load(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) load() error {
	data, err := os.ReadFile(dbPath)
	if os.IsNotExist(err) {
		db.users = []User{}
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &db.users)
}

func (db *Database) persist() error {
	data, err := json.MarshalIndent(db.users, "", "  ")
	if err != nil {
		return err
	}
	tmp := dbPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, dbPath)
}

// List returns all users, optionally filtered by a search string (name or email).
func (db *Database) List(search string) []User {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if search == "" {
		result := make([]User, len(db.users))
		copy(result, db.users)
		return result
	}

	lower := strings.ToLower(search)
	var result []User
	for _, u := range db.users {
		if strings.Contains(strings.ToLower(u.Name), lower) ||
			strings.Contains(strings.ToLower(u.Email), lower) {
			result = append(result, u)
		}
	}
	return result
}

// Create adds a new user and persists the database.
func (db *Database) Create(user User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.users = append(db.users, user)
	return db.persist()
}

// Update modifies name and/or email of the user with the given id.
// Returns false if no user with that id was found.
func (db *Database) Update(id, name, email string) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, u := range db.users {
		if u.ID == id {
			if name != "" {
				db.users[i].Name = name
			}
			if email != "" {
				db.users[i].Email = email
			}
			db.users[i].UpdatedAt = time.Now()
			return true, db.persist()
		}
	}
	return false, nil
}

// Delete removes the user with the given id.
// Returns false if no user with that id was found.
func (db *Database) Delete(id string) (bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, u := range db.users {
		if u.ID == id {
			db.users = append(db.users[:i], db.users[i+1:]...)
			return true, db.persist()
		}
	}
	return false, nil
}
