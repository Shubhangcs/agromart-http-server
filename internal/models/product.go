package models

import "time"

// --- DB Models ---

type Product struct {
	ID              string          `json:"id,omitempty"`
	UserID          string          `json:"user_id,omitempty"`
	BusinessID      string          `json:"business_id,omitempty"`
	CategoryID      string          `json:"category_id,omitempty"`
	SubCategoryID   string          `json:"sub_category_id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Quantity        float64         `json:"quantity,omitempty"`
	Unit            string          `json:"unit,omitempty"`
	Price           float64         `json:"price,omitempty"`
	MOQ             string          `json:"moq,omitempty"`
	Images          []ProductImages `json:"product_images,omitempty"`
	IsProductActive bool            `json:"is_product_active,omitempty"`
	CreatedAT       time.Time       `json:"created_at"`
	UpdatedAT       time.Time       `json:"updated_at"`
}

// ProductImages is the internal DB scan struct for product images.
type ProductImages struct {
	ID         string    `json:"id"`
	ImageIndex int       `json:"image_index"`
	ProductID  string    `json:"product_id"`
	Image      string    `json:"image"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

// ProductImage is the response-safe image struct.
type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Index     int       `json:"index"`
	Image     string    `json:"image"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

// ProductImageDetails is used when returning image id/index pairs.
type ProductImageDetails struct {
	ID    string `json:"id"`
	Index int    `json:"index"`
}

// --- Request DTOs ---

// CreateProductRequest is the payload for creating a product.
// swagger:model
type CreateProductRequest struct {
	BusinessID      string  `json:"business_id"       validate:"required" example:"biz-uuid-001"`
	CategoryID      string  `json:"category_id"       validate:"required" example:"cat-uuid-001"`
	SubCategoryID   string  `json:"sub_category_id"   validate:"required" example:"subcat-uuid-001"`
	Name            string  `json:"name"              validate:"required" example:"Wheat"`
	Description     string  `json:"description"       validate:"required" example:"High quality wheat"`
	Quantity        float64 `json:"quantity"                             example:"1000"`
	Unit            string  `json:"unit"              validate:"required" example:"kg"`
	Price           float64 `json:"price"                                example:"25.50"`
	MOQ             string  `json:"moq"               validate:"required" example:"100kg"`
	IsProductActive bool    `json:"is_product_active"                    example:"true"`
}

// UpdateProductRequest is the payload for updating a product.
// swagger:model
type UpdateProductRequest struct {
	Name          string  `json:"name"            example:"Wheat"`
	CategoryID    string  `json:"category_id"     example:"cat-uuid-001"`
	SubCategoryID string  `json:"sub_category_id" example:"subcat-uuid-001"`
	Quantity      float64 `json:"quantity"        example:"1000"`
	Price         float64 `json:"price"           example:"25.50"`
	Unit          string  `json:"unit"            example:"kg"`
	MOQ           string  `json:"moq"             example:"100kg"`
	Description   string  `json:"description"     example:"High quality wheat"`
}

// ChangeProductStatusRequest is the payload for toggling product active status.
// swagger:model
type ChangeProductStatusRequest struct {
	IsProductActive bool `json:"is_product_active" example:"true"`
}

// --- Response DTOs ---

// ProductResponse is the safe public representation of a product listing.
// swagger:model
type ProductResponse struct {
	ID              string          `json:"id"`
	BusinessID      string          `json:"business_id"`
	CategoryID      string          `json:"category_id"`
	CategoryName    string          `json:"category_name"`
	SubCategoryID   string          `json:"sub_category_id"`
	SubCategoryName string          `json:"sub_category_name"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Quantity        float64         `json:"quantity"`
	Unit            string          `json:"unit"`
	Price           float64         `json:"price"`
	MOQ             string          `json:"moq"`
	IsProductActive bool            `json:"is_product_active"`
	Images          []ProductImages `json:"images,omitempty"`
	CreatedAT       time.Time       `json:"created_at"`
	UpdatedAT       time.Time       `json:"updated_at"`
}

// ProductDetailsResponse is the enriched product detail response (replaces CompleteProduct).
// swagger:model
type ProductDetailsResponse struct {
	ID                     string          `json:"id"`
	UserID                 string          `json:"user_id"`
	BusinessID             string          `json:"business_id"`
	BusinessName           string          `json:"business_name"`
	BusinessEmail          string          `json:"business_email"`
	BusinessPhone          string          `json:"business_phone"`
	Address                string          `json:"address"`
	City                   string          `json:"city"`
	State                  string          `json:"state"`
	Pincode                string          `json:"pincode"`
	CategoryID             string          `json:"category_id"`
	CategoryName           string          `json:"category_name"`
	CategoryDescription    string          `json:"category_description"`
	SubCategoryID          string          `json:"sub_category_id"`
	SubCategoryName        string          `json:"sub_category_name"`
	SubCategoryDescription string          `json:"sub_category_description"`
	ProductName            string          `json:"product_name"`
	ProductDescription     string          `json:"product_description"`
	Quantity               float64         `json:"quantity"`
	Unit                   string          `json:"unit"`
	Price                  float64         `json:"price"`
	IsProductActive        bool            `json:"is_product_active"`
	MOQ                    string          `json:"moq"`
	CreatedAT              string          `json:"created_at"`
	Images                 []ProductImages `json:"product_images,omitempty"`
}
