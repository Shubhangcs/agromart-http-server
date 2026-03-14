package models

import "time"

// ProductRating represents a user's rating for a product.
type ProductRating struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name,omitempty"`
	Rating    float64   `json:"rating"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

// BusinessRating represents a user's rating for a business.
type BusinessRating struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"business_id"`
	UserID     string    `json:"user_id"`
	UserName   string    `json:"user_name,omitempty"`
	Rating     float64   `json:"rating"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}
