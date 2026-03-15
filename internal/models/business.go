package models

import "time"

// --- DB Models ---

type Business struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ProfileImage       *string   `json:"profile_image"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Pincode            string    `json:"pincode"`
	BusinessType       string    `json:"business_type"`
	IsBusinessVerified bool      `json:"is_business_verified"`
	IsBusinessTrusted  bool      `json:"is_business_trusted"`
	IsBusinessApproved bool      `json:"is_business_approved"`
	CreatedAT          time.Time `json:"created_at"`
	UpdatedAT          time.Time `json:"updated_at"`
}

type Social struct {
	ID        string    `json:"id"`
	Linkedin  *string   `json:"linkedin"`
	Instagram *string   `json:"instagram"`
	Telegram  *string   `json:"telegram"`
	Youtube   *string   `json:"youtube"`
	X         *string   `json:"x"`
	Facebook  *string   `json:"facebook"`
	Website   *string   `json:"website"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
}

type Legal struct {
	ID           string    `json:"id"`
	Aadhaar      *string   `json:"aadhaar"`
	Pan          *string   `json:"pan"`
	ExportImport *string   `json:"export_import"`
	MSME         *string   `json:"msme"`
	Fassi        *string   `json:"fassi"`
	GST          *string   `json:"gst"`
	CreatedAT    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
}

type BusinessApplication struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	RejectReason *string   `json:"reject_reason"`
	CreatedAT    time.Time `json:"created_at"`
}

type BusinessDetails struct {
	CoreBusinessDetails        Business            `json:"business_details"`
	BusinessSocialDetails      Social              `json:"social_details"`
	BusinessLegalDetails       Legal               `json:"legal_details"`
	BusinessApplicationDetails BusinessApplication `json:"business_application"`
}

// --- Request DTOs ---

// CreateBusinessRequest is the payload for creating a new business.
// swagger:model
type CreateBusinessRequest struct {
	UserID       string `json:"user_id"       validate:"required"       example:"user-uuid-001"`
	Name         string `json:"name"          validate:"required"       example:"Agro Traders Pvt Ltd"`
	Email        string `json:"email"         validate:"required,email" example:"business@example.com"`
	Phone        string `json:"phone"         validate:"required,phone" example:"9876543210"`
	Address      string `json:"address"       validate:"required"       example:"123 Market Street"`
	City         string `json:"city"          validate:"required"       example:"Pune"`
	State        string `json:"state"         validate:"required"       example:"Maharashtra"`
	Pincode      string `json:"pincode"       validate:"required"       example:"411001"`
	BusinessType string `json:"business_type" validate:"required"       example:"TRADER"`
}

// UpdateBusinessRequest is the payload for updating a business profile.
// swagger:model
type UpdateBusinessRequest struct {
	Name         string `json:"name"          example:"Agro Traders Pvt Ltd"`
	Email        string `json:"email"         validate:"omitempty,email"  example:"business@example.com"`
	Phone        string `json:"phone"         validate:"omitempty,phone"  example:"9876543210"`
	Address      string `json:"address"       example:"123 Market Street"`
	City         string `json:"city"          example:"Pune"`
	State        string `json:"state"         example:"Maharashtra"`
	Pincode      string `json:"pincode"       example:"411001"`
	BusinessType string `json:"business_type" example:"TRADER"`
}

// CreateSocialRequest is the payload for business social media links.
// swagger:model
type CreateSocialRequest struct {
	ID        string  `json:"id"        validate:"required" example:"biz-uuid-001"`
	Linkedin  *string `json:"linkedin"  example:"https://linkedin.com/company/agro"`
	Instagram *string `json:"instagram" example:"https://instagram.com/agro"`
	Telegram  *string `json:"telegram"  example:"https://t.me/agro"`
	Youtube   *string `json:"youtube"   example:"https://youtube.com/agro"`
	X         *string `json:"x"         example:"https://x.com/agro"`
	Facebook  *string `json:"facebook"  example:"https://facebook.com/agro"`
	Website   *string `json:"website"   example:"https://agromart.com"`
}

// CreateLegalRequest is the payload for business legal documents.
// swagger:model
type CreateLegalRequest struct {
	ID           string  `json:"id"            validate:"required"                example:"biz-uuid-001"`
	Aadhaar      *string `json:"aadhaar"       validate:"omitempty,aadhaar"       example:"123456789012"`
	Pan          *string `json:"pan"           validate:"omitempty,pan"           example:"ABCDE1234F"`
	ExportImport *string `json:"export_import" validate:"omitempty,export_import" example:"ABCD1234EF"`
	MSME         *string `json:"msme"          validate:"omitempty,msme"          example:"UDYAM-MH-01-0000001"`
	Fassi        *string `json:"fassi"         validate:"omitempty,fassi"         example:"12345678901234"`
	GST          *string `json:"gst"           validate:"omitempty,gst"           example:"29ABCDE1234F1Z5"`
}

// CreateApplicationRequest is the payload for submitting a business application.
// swagger:model
type CreateApplicationRequest struct {
	ID string `json:"id" validate:"required" example:"biz-uuid-001"`
}

// RejectApplicationRequest is the payload for rejecting a business application.
// swagger:model
type RejectApplicationRequest struct {
	RejectReason string `json:"reject_reason" validate:"required" example:"Documents incomplete"`
}

// UpdateBusinessStatusRequest is the payload for toggling business status flags.
// swagger:model
type UpdateBusinessStatusRequest struct {
	Status bool `json:"status" example:"true"`
}

// --- Response DTOs ---

// BusinessResponse is the safe public representation of a business.
// Business and its related structs (Social, Legal) have no sensitive fields
// and can be returned directly. This type is provided for explicit DTO usage.
// swagger:model
type BusinessResponse struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ProfileImage       *string   `json:"profile_image"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Phone              string    `json:"phone"`
	Address            string    `json:"address"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Pincode            string    `json:"pincode"`
	BusinessType       string    `json:"business_type"`
	IsBusinessVerified bool      `json:"is_business_verified"`
	IsBusinessTrusted  bool      `json:"is_business_trusted"`
	IsBusinessApproved bool      `json:"is_business_approved"`
	CreatedAT          time.Time `json:"created_at"`
	UpdatedAT          time.Time `json:"updated_at"`
}
