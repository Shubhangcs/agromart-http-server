package models

import "time"

type Product struct {
	ID              string          `json:"id,omitempty"`
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

type CompleteProduct struct {
	ID                     string          `json:"id"`
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

type ProductImages struct {
	ID         string    `json:"id"`
	ImageIndex int       `json:"image_index"`
	ProductID  string    `json:"product_id"`
	Image      string    `json:"image"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

type ProductImage struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Index     int       `json:"index"`
	Image     string    `json:"image"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type ProductImageDetails struct {
	ID    string `json:"id"`
	Index int    `json:"index"`
}
