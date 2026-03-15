package models

import "time"

// --- DB Models ---

// Message represents a single direct message between two users.
type Message struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAT  time.Time `json:"created_at"`
}

// --- Request DTOs ---

// SendMessageRequest is the payload for sending a direct message.
// swagger:model
type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" validate:"required" example:"user-uuid-002"`
	Content    string `json:"content"     validate:"required" example:"Hello, is this product available?"`
}
