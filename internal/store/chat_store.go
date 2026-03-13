package store

import (
	"database/sql"

	"github.com/shubhangcs/agromart-server/internal/models"
)

type PostgresChatStore struct {
	db *sql.DB
}

type ChatStore interface {
	// SaveMessage persists a new message and populates ID and CreatedAT.
	SaveMessage(*models.Message) error
	// GetChatHistory returns all messages between two users ordered oldest→newest.
	GetChatHistory(user1ID, user2ID string, limit, offset int) ([]models.Message, error)
	// MarkAsRead marks all unread messages sent to receiverID from senderID as read.
	MarkAsRead(senderID, receiverID string) error
}

func NewPostgresChatStore(db *sql.DB) *PostgresChatStore {
	return &PostgresChatStore{db: db}
}

func (cs *PostgresChatStore) SaveMessage(m *models.Message) error {
	query := `
	INSERT INTO messages (sender_id, receiver_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, is_read, created_at
	`
	return cs.db.QueryRow(query, m.SenderID, m.ReceiverID, m.Content).
		Scan(&m.ID, &m.IsRead, &m.CreatedAT)
}

// GetChatHistory returns the conversation between two users, oldest first,
// using LEAST/GREATEST so the query works regardless of sender/receiver order.
func (cs *PostgresChatStore) GetChatHistory(user1ID, user2ID string, limit, offset int) ([]models.Message, error) {
	query := `
	SELECT id, sender_id, receiver_id, content, is_read, created_at
	FROM messages
	WHERE (sender_id = $1 AND receiver_id = $2)
	   OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY created_at ASC
	LIMIT $3 OFFSET $4
	`
	rows, err := cs.db.Query(query, user1ID, user2ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err = rows.Scan(
			&msg.ID, &msg.SenderID, &msg.ReceiverID,
			&msg.Content, &msg.IsRead, &msg.CreatedAT,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// MarkAsRead marks all messages from senderID to receiverID as read.
func (cs *PostgresChatStore) MarkAsRead(senderID, receiverID string) error {
	_, err := cs.db.Exec(
		`UPDATE messages SET is_read = TRUE
		 WHERE sender_id = $1 AND receiver_id = $2 AND is_read = FALSE`,
		senderID, receiverID,
	)
	return err
}
