package models

import "time"

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

