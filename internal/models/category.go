package models

import "time"

// --- DB Models ---

type Category struct {
	ID            string    `json:"id"`
	CategoryImage *string   `json:"category_image"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

type SubCategory struct {
	ID               string    `json:"id"`
	CategoryID       string    `json:"category_id"`
	SubCategoryImage *string   `json:"sub_category_image"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	CreatedAT        time.Time `json:"created_at"`
	UpdatedAT        time.Time `json:"updated_at"`
}

// --- Request DTOs ---

// CreateCategoryRequest is the payload for creating a product category.
// swagger:model
type CreateCategoryRequest struct {
	Name        string `json:"name"        validate:"required" example:"Grains"`
	Description string `json:"description" validate:"required" example:"All types of grains and cereals"`
}

// UpdateCategoryRequest is the payload for updating a category.
// swagger:model
type UpdateCategoryRequest struct {
	Name        string `json:"name"        example:"Grains"`
	Description string `json:"description" example:"All types of grains and cereals"`
}

// CreateSubCategoryRequest is the payload for creating a sub-category.
// swagger:model
type CreateSubCategoryRequest struct {
	CategoryID  string `json:"category_id"  validate:"required" example:"cat-uuid-001"`
	Name        string `json:"name"         validate:"required" example:"Wheat"`
	Description string `json:"description"  validate:"required" example:"All varieties of wheat"`
}

// UpdateSubCategoryRequest is the payload for updating a sub-category.
// swagger:model
type UpdateSubCategoryRequest struct {
	Name        string `json:"name"        example:"Wheat"`
	Description string `json:"description" example:"All varieties of wheat"`
}
