package models

import "time"

// BusinessReview represents a user's written review for a business.
type BusinessReview struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"business_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"`
	Review     string    `json:"review"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

// ProductReview represents a user's written review for a product.
type ProductReview struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Review    string    `json:"review"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}
