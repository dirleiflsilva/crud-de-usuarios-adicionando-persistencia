package database

import (
	"database/sql"
	"errors"
	"fmt"
	"usuarios-crud/models"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var ErrUserNotFound = errors.New("user not found")

// InitDB abre a conexao e garante que a tabela exista antes da API iniciar.
func InitDB(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Ping valida se a conexao foi aberta corretamente.
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db

	// Se a tabela ainda nao existir, ela e criada aqui.
	err = createUserTable()
	if err != nil {
		return err
	}

	fmt.Println("Database initialized successfully")
	return nil
}

func CreateUser(user *models.User) error {
	query := `
	INSERT INTO users (id, first_name, last_name, biography)
	VALUES (?, ?, ?, ?)
	`

	_, err := DB.Exec(query, user.ID, user.FirstName, user.LastName, user.Biography)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, first_name, last_name, biography FROM users WHERE id = ?`

	user := &models.User{}
	// Scan copia as colunas retornadas pelo SELECT para a struct.
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Biography)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	query := `SELECT id, first_name, last_name, biography FROM users`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		// Cada linha lida do banco vira um item do slice users.
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Biography)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

func UpdateUser(id string, user *models.User) error {
	query := `
	UPDATE users
	SET first_name = ?, last_name = ?, biography = ?
	WHERE id = ?
	`

	result, err := DB.Exec(query, user.FirstName, user.LastName, user.Biography, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func DeleteUser(id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// createUserTable executa o schema inicial da aplicacao.
func createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		biography TEXT NOT NULL
	)
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}

// CloseDB encerra a conexao ao finalizar a aplicacao.
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
