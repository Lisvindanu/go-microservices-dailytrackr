package models

import (
	"database/sql"
	"time"
)

// User represents the user model
type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Hidden from JSON
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = ?
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE username = ?
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&count)
	return count > 0, err
}

// UsernameExists checks if a username already exists
func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = ?"
	err := r.db.QueryRow(query, username).Scan(&count)
	return count > 0, err
}

// Update updates user information
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?
	`

	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
