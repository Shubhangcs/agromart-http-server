package models

import "time"

// --- Request DTOs ---

// AddToWishlistRequest is the request payload for adding a product to the wishlist.
// swagger:model
type AddToWishlistRequest struct {
	ProductID string `json:"product_id" validate:"required" example:"prod-uuid-001"`
}

// --- Response DTOs ---

// WishlistItem represents a single product saved in a user's wishlist,
// joined with full product details for display.
// swagger:model
type WishlistItem struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	Unit         string    `json:"unit"`
	MOQ          string    `json:"moq"`
	BusinessID   string    `json:"business_id"`
	ProductImage *string   `json:"product_image"`
	CreatedAT    time.Time `json:"created_at"`
}
