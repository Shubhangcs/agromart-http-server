package store

import "database/sql"

type PushStore interface {
	Upsert(userID, token, platform string) error
	GetTokens(userID string) ([]string, error)
	Delete(token string) error
}

type PostgresPushStore struct{ db *sql.DB }

func NewPostgresPushStore(db *sql.DB) *PostgresPushStore { return &PostgresPushStore{db: db} }

// Upsert registers a device token for a user; a token that moves to another account is re-owned.
func (s *PostgresPushStore) Upsert(userID, token, platform string) error {
	_, err := s.db.Exec(`
		INSERT INTO push_tokens (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform, updated_at = NOW()`,
		userID, token, platform)
	return err
}

func (s *PostgresPushStore) GetTokens(userID string) ([]string, error) {
	rows, err := s.db.Query(`SELECT token FROM push_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tokens []string
	for rows.Next() {
		var t string
		if err = rows.Scan(&t); err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, rows.Err()
}

func (s *PostgresPushStore) Delete(token string) error {
	_, err := s.db.Exec(`DELETE FROM push_tokens WHERE token = $1`, token)
	return err
}
