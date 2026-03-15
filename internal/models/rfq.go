package models

import "time"

// --- DB Models ---

// RFQ is the core database model for a Request for Quotation.
type RFQ struct {
	ID            string    `json:"id"`
	BusinessID    string    `json:"business_id"`
	CategoryID    string    `json:"category_id"`
	SubCategoryID string    `json:"sub_category_id"`
	ProductName   string    `json:"product_name"`
	Quantity      float64   `json:"quantity"`
	Unit          string    `json:"unit"`
	Price         float64   `json:"price"`
	IsRFQActive   bool      `json:"is_rfq_active"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

// --- Request DTOs ---

// CreateRFQRequest is the payload for creating a Request for Quotation.
// swagger:model
type CreateRFQRequest struct {
	BusinessID    string  `json:"business_id"     validate:"required" example:"biz-uuid-001"`
	CategoryID    string  `json:"category_id"     validate:"required" example:"cat-uuid-001"`
	SubCategoryID string  `json:"sub_category_id" validate:"required" example:"subcat-uuid-001"`
	ProductName   string  `json:"product_name"    validate:"required" example:"Wheat"`
	Quantity      float64 `json:"quantity"                            example:"500"`
	Unit          string  `json:"unit"            validate:"required" example:"kg"`
	Price         float64 `json:"price"                               example:"1200.50"`
	IsRFQActive   bool    `json:"is_rfq_active"                       example:"true"`
}

// UpdateRFQRequest is the payload for updating an RFQ.
// swagger:model
type UpdateRFQRequest struct {
	CategoryID    string  `json:"category_id"     example:"cat-uuid-001"`
	SubCategoryID string  `json:"sub_category_id" example:"subcat-uuid-001"`
	ProductName   string  `json:"product_name"    example:"Wheat"`
	Quantity      float64 `json:"quantity"        example:"500"`
	Price         float64 `json:"price"           example:"1200.50"`
	Unit          string  `json:"unit"            example:"kg"`
}

// ActivateRFQRequest is the payload for toggling an RFQ's active status.
// swagger:model
type ActivateRFQRequest struct {
	IsRFQActive bool `json:"is_rfq_active" example:"true"`
}

// --- Response DTOs ---

// RFQResponse is the enriched RFQ response with joined business and category details.
// swagger:model
type RFQResponse struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	BusinessID             string    `json:"business_id"`
	BusinessName           string    `json:"business_name"`
	BusinessPhone          string    `json:"business_phone"`
	BusinessEmail          string    `json:"business_email"`
	Address                string    `json:"address"`
	City                   string    `json:"city"`
	State                  string    `json:"state"`
	CategoryID             string    `json:"category_id,omitempty"`
	SubCategoryID          string    `json:"sub_category_id,omitempty"`
	CategoryName           string    `json:"category_name,omitempty"`
	SubCategoryName        string    `json:"sub_category_name,omitempty"`
	CategoryDescription    string    `json:"category_description,omitempty"`
	SubCategoryDescription string    `json:"sub_category_description,omitempty"`
	ProductName            string    `json:"product_name"`
	Quantity               float64   `json:"quantity"`
	Unit                   string    `json:"unit"`
	Price                  float64   `json:"price"`
	IsRFQActive            bool      `json:"is_rfq_active"`
	CreatedAT              time.Time `json:"created_at"`
	UpdatedAT              time.Time `json:"updated_at"`
}
