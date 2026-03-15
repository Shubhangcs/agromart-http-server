package models

import "time"

// --- DB Models (also serve as Response DTOs - UserName is joined) ---

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

// --- Request DTOs ---

// RateProductRequest is the payload for rating a product.
// swagger:model
type RateProductRequest struct {
	ProductID string  `json:"product_id" validate:"required" example:"prod-uuid-001"`
	UserID    string  `json:"user_id"    validate:"required" example:"user-uuid-001"`
	Rating    float64 `json:"rating"                         example:"4.5"`
}

// RateBusinessRequest is the payload for rating a business.
// swagger:model
type RateBusinessRequest struct {
	BusinessID string  `json:"business_id" validate:"required" example:"biz-uuid-001"`
	UserID     string  `json:"user_id"     validate:"required" example:"user-uuid-001"`
	Rating     float64 `json:"rating"                          example:"4.5"`
}
