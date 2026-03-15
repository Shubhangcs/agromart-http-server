package models

import "time"

// --- DB Models (also serve as Response DTOs - UserName is joined) ---

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

// --- Request DTOs ---

// CreateBusinessReviewRequest is the payload for creating a business review.
// swagger:model
type CreateBusinessReviewRequest struct {
	BusinessID string `json:"business_id" validate:"required" example:"biz-uuid-001"`
	UserID     string `json:"user_id"     validate:"required" example:"user-uuid-001"`
	Review     string `json:"review"      validate:"required" example:"Excellent service!"`
}

// UpdateReviewRequest is the payload for updating a business or product review.
// swagger:model
type UpdateReviewRequest struct {
	Review string `json:"review" validate:"required" example:"Updated review text"`
}

// CreateProductReviewRequest is the payload for creating a product review.
// swagger:model
type CreateProductReviewRequest struct {
	ProductID string `json:"product_id" validate:"required" example:"prod-uuid-001"`
	UserID    string `json:"user_id"    validate:"required" example:"user-uuid-001"`
	Review    string `json:"review"     validate:"required" example:"Great product quality!"`
}
