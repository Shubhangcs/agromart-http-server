package models

// --- User / Admin DTOs ---

// CreateAdminRequest is the payload for creating a new admin.
// swagger:model
type CreateAdminRequest struct {
	FirstName string  `json:"first_name" validate:"required"       example:"John"`
	LastName  *string `json:"last_name"                            example:"Doe"`
	Email     string  `json:"email"      validate:"required,email" example:"admin@example.com"`
	Phone     string  `json:"phone"      validate:"required,phone" example:"9876543210"`
	Password  string  `json:"password"   validate:"required"       example:"secret123"`
}

// UpdateUserDetailsRequest is the payload for updating user/admin profile details.
// swagger:model
type UpdateUserDetailsRequest struct {
	FirstName string  `json:"first_name"                              example:"John"`
	LastName  *string `json:"last_name"                               example:"Doe"`
	Email     string  `json:"email"      validate:"omitempty,email"   example:"user@example.com"`
	Phone     string  `json:"phone"      validate:"omitempty,phone"   example:"9876543210"`
}

// UpdatePasswordRequest is the payload for changing a password.
// swagger:model
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required" example:"oldpass123"`
	NewPassword string `json:"new_password" validate:"required" example:"newpass456"`
}

// LoginRequest is the payload for login endpoints.
// swagger:model
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required"       example:"secret123"`
}

// BlockUserRequest is the payload for blocking/unblocking a user.
// swagger:model
type BlockUserRequest struct {
	IsUserBlocked bool `json:"is_user_blocked" example:"true"`
}

// --- Business DTOs ---

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
	ID           string  `json:"id"            validate:"required"                  example:"biz-uuid-001"`
	Aadhaar      *string `json:"aadhaar"       validate:"omitempty,aadhaar"         example:"123456789012"`
	Pan          *string `json:"pan"           validate:"omitempty,pan"             example:"ABCDE1234F"`
	ExportImport *string `json:"export_import" validate:"omitempty,export_import"   example:"ABCD1234EF"`
	MSME         *string `json:"msme"          validate:"omitempty,msme"            example:"UDYAM-MH-01-0000001"`
	Fassi        *string `json:"fassi"         validate:"omitempty,fassi"           example:"12345678901234"`
	GST          *string `json:"gst"           validate:"omitempty,gst"             example:"29ABCDE1234F1Z5"`
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

// RateBusinessRequest is the payload for rating a business.
// swagger:model
type RateBusinessRequest struct {
	BusinessID string  `json:"business_id" validate:"required" example:"biz-uuid-001"`
	UserID     string  `json:"user_id"     validate:"required" example:"user-uuid-001"`
	Rating     float64 `json:"rating"                          example:"4.5"`
}

// --- Product DTOs ---

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

// --- RFQ DTOs ---

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

// --- Category DTOs ---

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

// --- Follower DTOs ---

// FollowRequest is the payload for follow/unfollow actions.
// swagger:model
type FollowRequest struct {
	UserID     string `json:"user_id"     validate:"required" example:"user-uuid-001"`
	BusinessID string `json:"business_id" validate:"required" example:"biz-uuid-001"`
}

// --- Rating / Review DTOs ---

// RateProductRequest is the payload for rating a product.
// swagger:model
type RateProductRequest struct {
	ProductID string  `json:"product_id" validate:"required" example:"prod-uuid-001"`
	UserID    string  `json:"user_id"    validate:"required" example:"user-uuid-001"`
	Rating    float64 `json:"rating"                         example:"4.5"`
}

// CreateBusinessReviewRequest is the payload for creating a business review.
// swagger:model
type CreateBusinessReviewRequest struct {
	BusinessID string `json:"business_id" validate:"required" example:"biz-uuid-001"`
	UserID     string `json:"user_id"     validate:"required" example:"user-uuid-001"`
	Review     string `json:"review"      validate:"required" example:"Excellent service!"`
}

// UpdateReviewRequest is the payload for updating a business or product review.
// swagger:model
type UpdateReviewRequest struct {
	Review string `json:"review" validate:"required" example:"Updated review text"`
}

// CreateProductReviewRequest is the payload for creating a product review.
// swagger:model
type CreateProductReviewRequest struct {
	ProductID string `json:"product_id" validate:"required" example:"prod-uuid-001"`
	UserID    string `json:"user_id"    validate:"required" example:"user-uuid-001"`
	Review    string `json:"review"     validate:"required" example:"Great product quality!"`
}

// --- Chat DTOs ---

// SendMessageRequest is the payload for sending a direct message.
// swagger:model
type SendMessageRequest struct {
	ReceiverID string `json:"receiver_id" validate:"required" example:"user-uuid-002"`
	Content    string `json:"content"     validate:"required" example:"Hello, is this product available?"`
}
