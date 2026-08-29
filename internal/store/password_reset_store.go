package store

import (
	"database/sql"
	"time"
)

type PasswordReset struct {
	ID        string
	Email     string
	CodeHash  string
	Attempts  int
	ExpiresAt time.Time
}

type PasswordResetStore interface {
	Create(email, codeHash string, expiresAt time.Time) error
	GetActive(email string) (*PasswordReset, error)
	IncrementAttempts(id string) error
	MarkUsed(id string) error
}

type PostgresPasswordResetStore struct{ db *sql.DB }

func NewPostgresPasswordResetStore(db *sql.DB) *PostgresPasswordResetStore {
	return &PostgresPasswordResetStore{db: db}
}

// Create invalidates any earlier codes for the email and stores a new one.
func (s *PostgresPasswordResetStore) Create(email, codeHash string, expiresAt time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec(`UPDATE password_resets SET used = TRUE WHERE email = $1 AND used = FALSE`, email); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO password_resets (email, code_hash, expires_at) VALUES ($1, $2, $3)`, email, codeHash, expiresAt); err != nil {
		return err
	}
	return tx.Commit()
}

// GetActive returns the newest unused, unexpired code for the email (sql.ErrNoRows if none).
func (s *PostgresPasswordResetStore) GetActive(email string) (*PasswordReset, error) {
	var pr PasswordReset
	err := s.db.QueryRow(`
		SELECT id, email, code_hash, attempts, expires_at
		FROM password_resets
		WHERE email = $1 AND used = FALSE AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1`, email).Scan(&pr.ID, &pr.Email, &pr.CodeHash, &pr.Attempts, &pr.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

func (s *PostgresPasswordResetStore) IncrementAttempts(id string) error {
	_, err := s.db.Exec(`UPDATE password_resets SET attempts = attempts + 1 WHERE id = $1`, id)
	return err
}

func (s *PostgresPasswordResetStore) MarkUsed(id string) error {
	_, err := s.db.Exec(`UPDATE password_resets SET used = TRUE WHERE id = $1`, id)
	return err
}
