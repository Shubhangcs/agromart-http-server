package models

import "time"

type RFQ struct {
	ID                     string    `json:"id"`
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
