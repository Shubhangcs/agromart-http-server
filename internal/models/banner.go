package models

import "time"

type Banner struct {
	ID        string    `json:"id"`
	Title     *string   `json:"title"`
	Image     *string   `json:"image"`
	TargetURL *string   `json:"target_url"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type UpsertBannerRequest struct {
	Title     *string `json:"title"      example:"Monsoon cashew offers"`
	TargetURL *string `json:"target_url" example:"southcanra://pages/productsByCategory?category_id=..."`
	SortOrder int     `json:"sort_order" example:"1"`
	IsActive  *bool   `json:"is_active"  example:"true"`
}
