package models

import "time"

// --- DB Model ---

// Lead is the database representation of a product enquiry.
// swagger:model
type Lead struct {
	LeadID             string    `json:"lead_id"`
	EnquirerBusinessID string    `json:"enquirer_business_id"`
	EnquireToID        string    `json:"enquire_to_id"`
	ProductID          string    `json:"product_id"`
	EnquiryMessage     string    `json:"enquiry_message"`
	OrderQuantity      float64   `json:"order_quantity"`
	ExpectedPrice      float64   `json:"expected_price"`
	CreatedAT          time.Time `json:"created_at"`
}

// --- Request DTOs ---

// CreateLeadRequest is the payload for submitting a product enquiry.
// swagger:model
type CreateLeadRequest struct {
	EnquirerBusinessID string  `json:"enquirer_business_id" validate:"required"`
	EnquireToID        string  `json:"enquire_to_id"        validate:"required"`
	ProductID          string  `json:"product_id"           validate:"required"`
	EnquiryMessage     string  `json:"enquiry_message"      validate:"required"`
	OrderQuantity      float64 `json:"order_quantity"`
	ExpectedPrice      float64 `json:"expected_price"`
}

// --- Response DTOs ---

// LeadSentResponse is returned to the business that sent the enquiry.
// It includes the enquired product details and the complete seller business info.
// swagger:model
type LeadSentResponse struct {
	LeadID                    string    `json:"lead_id"`
	ProductID                 string    `json:"product_id"`
	ProductName               string    `json:"product_name"`
	ProductImage              *string   `json:"product_image"`
	SellerBusinessID          string    `json:"seller_business_id"`
	SellerBusinessName        string    `json:"seller_business_name"`
	SellerBusinessEmail       string    `json:"seller_business_email"`
	SellerBusinessPhone       string    `json:"seller_business_phone"`
	SellerBusinessProfileImage *string  `json:"seller_business_profile_image"`
	SellerAddress             string    `json:"seller_address"`
	SellerCity                string    `json:"seller_city"`
	SellerState               string    `json:"seller_state"`
	SellerPincode             string    `json:"seller_pincode"`
	EnquiryMessage            string    `json:"enquiry_message"`
	OrderQuantity             float64   `json:"order_quantity"`
	ExpectedPrice             float64   `json:"expected_price"`
	CreatedAT                 time.Time `json:"created_at"`
}

// LeadReceivedResponse is returned to the business that received the enquiry.
// It includes the enquired product details and the complete enquirer business info.
// swagger:model
type LeadReceivedResponse struct {
	LeadID                       string    `json:"lead_id"`
	ProductID                    string    `json:"product_id"`
	ProductName                  string    `json:"product_name"`
	ProductImage                 *string   `json:"product_image"`
	EnquirerBusinessID           string    `json:"enquirer_business_id"`
	EnquirerBusinessName         string    `json:"enquirer_business_name"`
	EnquirerBusinessEmail        string    `json:"enquirer_business_email"`
	EnquirerBusinessPhone        string    `json:"enquirer_business_phone"`
	EnquirerBusinessProfileImage *string   `json:"enquirer_business_profile_image"`
	EnquirerAddress              string    `json:"enquirer_address"`
	EnquirerCity                 string    `json:"enquirer_city"`
	EnquirerState                string    `json:"enquirer_state"`
	EnquirerPincode              string    `json:"enquirer_pincode"`
	EnquiryMessage               string    `json:"enquiry_message"`
	OrderQuantity                float64   `json:"order_quantity"`
	ExpectedPrice                float64   `json:"expected_price"`
	CreatedAT                    time.Time `json:"created_at"`
}
