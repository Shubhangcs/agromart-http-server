package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// businessRequest represents the business creation/update payload
type businessRequest struct {
	UserID       string `json:"user_id"       example:"user-uuid-001"`
	Name         string `json:"name"          example:"Agro Traders Pvt Ltd"`
	Email        string `json:"email"         example:"business@example.com"`
	Phone        string `json:"phone"         example:"9876543210"`
	Address      string `json:"address"       example:"123 Market Street"`
	City         string `json:"city"          example:"Pune"`
	State        string `json:"state"         example:"Maharashtra"`
	Pincode      string `json:"pincode"       example:"411001"`
	BusinessType string `json:"business_type" example:"TRADER"`
}

// socialRequest represents the business social media links payload
type socialRequest struct {
	ID        string  `json:"id"        example:"biz-uuid-001"`
	Linkedin  *string `json:"linkedin"  example:"https://linkedin.com/company/agro"`
	Instagram *string `json:"instagram" example:"https://instagram.com/agro"`
	Telegram  *string `json:"telegram"  example:"https://t.me/agro"`
	Youtube   *string `json:"youtube"   example:"https://youtube.com/agro"`
	X         *string `json:"x"         example:"https://x.com/agro"`
	Facebook  *string `json:"facebook"  example:"https://facebook.com/agro"`
	Website   *string `json:"website"   example:"https://agromart.com"`
}

// legalRequest represents the business legal documents payload
type legalRequest struct {
	ID           string  `json:"id"            example:"biz-uuid-001"`
	Aadhaar      *string `json:"aadhaar"       example:"1234-5678-9012"`
	Pan          *string `json:"pan"           example:"ABCDE1234F"`
	ExportImport *string `json:"export_import" example:"IEC123456"`
	MSME         *string `json:"msme"          example:"MSME123456"`
	Fassi        *string `json:"fassi"         example:"FASSI123456"`
	GST          *string `json:"gst"           example:"29ABCDE1234F1Z5"`
}

// businessApplicationRequest represents the business application payload
type businessApplicationRequest struct {
	ID           string `json:"id"            example:"biz-uuid-001"`
	Status       string `json:"status"        example:"APPLIED"`
	RejectReason string `json:"reject_reason" example:"Documents incomplete"`
}

type businessRatingRequest struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"business_id"`
	UserID     string    `json:"user_id"`
	Rating     float64   `json:"rating"`
	CreatedAT  time.Time `json:"created_at"`
	UpdatedAT  time.Time `json:"updated_at"`
}

type BusinessHandler struct {
	businessStore store.BusinessStore
	logger        *log.Logger
}

func NewBusinessHandler(businessStore store.BusinessStore, logger *log.Logger) *BusinessHandler {
	return &BusinessHandler{
		businessStore: businessStore,
		logger:        logger,
	}
}

func (bh *BusinessHandler) validateCreateBusinessRequest(req *businessRequest) error {
	if req.UserID == "" {
		return errors.New("invalid request user id is required")
	}
	if req.Name == "" {
		return errors.New("invalid request business name is required")
	}
	if !emailRegx.Match([]byte(req.Email)) {
		return errors.New("invalid request enter proper email id")
	}
	if !phoneRegx.Match([]byte(req.Phone)) {
		return errors.New("invalid request enter valid phone number")
	}
	if req.Address == "" {
		return errors.New("invalid request address is required")
	}
	if req.City == "" {
		return errors.New("invalid request city is required")
	}
	if req.State == "" {
		return errors.New("invalid request state is required")
	}
	if req.Pincode == "" {
		return errors.New("invalid request pincode is required")
	}
	if req.BusinessType == "" {
		return errors.New("invalid request business type is required")
	}
	return nil
}

func (bh *BusinessHandler) validateBusinessRatingRequest(req *businessRatingRequest) error {
	if req.BusinessID == "" {
		return errors.New("invalid request business id is required")
	}

	if req.UserID == "" {
		return errors.New("invalid request user id is required")
	}

	if req.Rating == 0.0 {
		return errors.New("invalid request rating is required")
	}

	return nil
}

// HandleCreateBusiness godoc
// @Summary      Create a business
// @Description  Creates a new business profile for a user
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body businessRequest true "Business creation payload"
// @Success      201 {object} map[string]interface{} "business created successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/create [post]
func (bh *BusinessHandler) HandleCreateBusiness(w http.ResponseWriter, r *http.Request) {
	var req businessRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: create business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = bh.validateCreateBusinessRequest(&req)
	if err != nil {
		bh.logger.Printf("ERROR: create business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	business := &store.Business{
		UserID:       req.UserID,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		BusinessType: req.BusinessType,
	}
	err = bh.businessStore.CreateBusiness(business)
	if err != nil {
		bh.logger.Printf("ERROR: create business: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "business created successfully", "business_id": business.ID})
}

// HandleCreateSocial godoc
// @Summary      Add business social links
// @Description  Creates social media links for the given business
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body socialRequest true "Social links payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing business ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/social/create [post]
func (bh *BusinessHandler) HandleCreateSocial(w http.ResponseWriter, r *http.Request) {
	var req socialRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: create social: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if req.ID == "" {
		bh.logger.Println("ERROR: create social: business id is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request business id is required"})
		return
	}
	social := &store.Social{
		ID:        req.ID,
		Linkedin:  req.Linkedin,
		Instagram: req.Instagram,
		Facebook:  req.Facebook,
		Youtube:   req.Youtube,
		Telegram:  req.Telegram,
		X:         req.X,
		Website:   req.Website,
	}
	err = bh.businessStore.CreateSocial(social)
	if err != nil {
		bh.logger.Printf("ERROR: create social: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "business socials created successfully"})
}

// HandleCreateLegal godoc
// @Summary      Add business legal documents
// @Description  Creates legal document records for the given business
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body legalRequest true "Legal documents payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing business ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/legal/create [post]
func (bh *BusinessHandler) HandleCreateLegal(w http.ResponseWriter, r *http.Request) {
	var req legalRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: create legal: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if req.ID == "" {
		bh.logger.Println("ERROR: create legal: business id is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request business id is required"})
		return
	}
	legal := &store.Legal{
		ID:           req.ID,
		Aadhaar:      req.Aadhaar,
		Pan:          req.Pan,
		ExportImport: req.ExportImport,
		MSME:         req.MSME,
		Fassi:        req.Fassi,
		GST:          req.GST,
	}
	err = bh.businessStore.CreateLegal(legal)
	if err != nil {
		bh.logger.Printf("ERROR: create legal: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "business legals created successfully"})
}

// HandleCreateBusinessApplication godoc
// @Summary      Submit a business application
// @Description  Submits an application for business approval with status set to APPLIED
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body businessApplicationRequest true "Business application payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing business ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/application/create [post]
func (bh *BusinessHandler) HandleCreateBusinessApplication(w http.ResponseWriter, r *http.Request) {
	var req businessApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: create business application: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if req.ID == "" {
		bh.logger.Println("ERROR: create business application: business id is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request business id is required"})
		return
	}
	businessApplication := &store.BusinessApplication{
		ID:     req.ID,
		Status: "APPLIED",
	}
	err = bh.businessStore.CreateBusinessApplication(businessApplication)
	if err != nil {
		bh.logger.Printf("ERROR: create business application: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "business application created successfully"})
}

// HandleUpdateBusiness godoc
// @Summary      Update business details
// @Description  Updates the profile details of the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int             true "Business ID"
// @Param        body body businessRequest true "Business update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateBusiness(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req businessRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	business := &store.Business{
		ID:           id,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		BusinessType: req.BusinessType,
	}
	err = bh.businessStore.UpdateBusiness(business)
	if err != nil {
		bh.logger.Printf("ERROR: update business: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business profile updated successfully"})
}

// HandleUpdateSocials godoc
// @Summary      Update business social links
// @Description  Updates the social media links for the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int           true "Business ID"
// @Param        body body socialRequest true "Social links update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/social/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateSocials(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update socials: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req socialRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update socials: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	social := &store.Social{
		ID:        id,
		Linkedin:  req.Linkedin,
		Instagram: req.Instagram,
		Youtube:   req.Youtube,
		Telegram:  req.Telegram,
		X:         req.X,
		Facebook:  req.Facebook,
		Website:   req.Website,
	}
	err = bh.businessStore.UpdateSocial(social)
	if err != nil {
		bh.logger.Printf("ERROR: update socials: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business social details updated successfully"})
}

// HandleUpdateLegals godoc
// @Summary      Update business legal documents
// @Description  Updates the legal document records for the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int          true "Business ID"
// @Param        body body legalRequest true "Legal documents update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/legal/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateLegals(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update legals: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req legalRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update legals: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	legal := &store.Legal{
		ID:           id,
		Aadhaar:      req.Aadhaar,
		Pan:          req.Pan,
		ExportImport: req.ExportImport,
		MSME:         req.MSME,
		Fassi:        req.Fassi,
		GST:          req.GST,
	}
	err = bh.businessStore.UpdateLegal(legal)
	if err != nil {
		bh.logger.Printf("ERROR: update legals: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business legal details updated successfully"})
}

// HandleAcceptBusinessApplication godoc
// @Summary      Accept a business application
// @Description  Marks the business application for the given ID as accepted
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/application/accept/{id} [put]
func (bh *BusinessHandler) HandleAcceptBusinessApplication(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: accept business application: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = bh.businessStore.AcceptBusinessApplication(id)
	if err != nil {
		bh.logger.Printf("ERROR: accept business application: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business application accepted successfully"})
}

// HandleRejectBusinessApplication godoc
// @Summary      Reject a business application
// @Description  Marks the business application for the given ID as rejected with a reason
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int                        true "Business ID"
// @Param        body body businessApplicationRequest true "Rejection reason payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID, payload, or missing reject reason"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/application/reject/{id} [put]
func (bh *BusinessHandler) HandleRejectBusinessApplication(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: reject business application: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req businessApplicationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: reject business application: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if req.RejectReason == "" {
		bh.logger.Println("ERROR: reject business application: reject reason is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request reject reason is required"})
		return
	}
	businessApplication := &store.BusinessApplication{
		ID:           id,
		Status:       "REJECTED",
		RejectReason: &req.RejectReason,
	}
	err = bh.businessStore.RejectBusinessApplication(businessApplication)
	if err != nil {
		bh.logger.Printf("ERROR: reject business application: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business application rejected successfully"})
}

// HandleGetCompleteBusinessDetails godoc
// @Summary      Get complete business details
// @Description  Returns full business profile including social, legal, and application details
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "complete business details"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/get/complete/{id} [get]
func (bh *BusinessHandler) HandleGetCompleteBusinessDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get complete business details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetCompleteBusinessDetails(id)
	if err != nil {
		bh.logger.Printf("ERROR: get complete business details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business details fetched successfully", "details": res})
}

// HandleGetBusinessDetails godoc
// @Summary      Get business details
// @Description  Returns the basic profile of the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "business details"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/get/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get business details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetBusiness(id)
	if err != nil {
		bh.logger.Printf("ERROR: get business details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business details fetched successfully", "details": res})
}

// HandleGetSocialDetails godoc
// @Summary      Get business social details
// @Description  Returns the social media links for the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "social details"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/social/get/{id} [get]
func (bh *BusinessHandler) HandleGetSocialDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get business social details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetSocial(id)
	if err != nil {
		bh.logger.Printf("ERROR: get business social details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business social details fetched successfully", "details": res})
}

// HandleGetLegalDetails godoc
// @Summary      Get business legal details
// @Description  Returns the legal documents for the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "legal details"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/legal/get/{id} [get]
func (bh *BusinessHandler) HandleGetLegalDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get business legal details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetLegal(id)
	if err != nil {
		bh.logger.Printf("ERROR: get business legal details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business legal details fetched successfully", "details": res})
}

// HandleGetBusinessApplicationDetails godoc
// @Summary      Get business application details
// @Description  Returns the application status and details for the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "application details"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/application/get/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessApplicationDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get business application details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetBusinessApplication(id)
	if err != nil {
		bh.logger.Printf("ERROR: get business application details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business application details fetched successfully", "details": res})
}

// HandleUpdateVerifyBusinessStatus godoc
// @Summary      Update business verified status
// @Description  Sets the verified status of the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int                   true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/status/verify/{id} [put]
func (bh *BusinessHandler) HandleUpdateVerifyBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update verify business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req struct {
		Status bool `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update verify business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = bh.businessStore.UpdateVerifyBusinessStatus(id, req.Status)
	if err != nil {
		bh.logger.Printf("ERROR: update verify business status: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business verified status updated successfully"})
}

// HandleUpdateTrustBusinessStatus godoc
// @Summary      Update business trusted status
// @Description  Sets the trusted status of the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int                   true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/status/trust/{id} [put]
func (bh *BusinessHandler) HandleUpdateTrustBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update trust business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req struct {
		Status bool `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update trust business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = bh.businessStore.UpdateTrustBusinessStatus(id, req.Status)
	if err != nil {
		bh.logger.Printf("ERROR: update trust business status: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business trusted status updated successfully"})
}

// HandleUpdateBlockBusinessStatus godoc
// @Summary      Update business blocked status
// @Description  Sets the blocked status of the business with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path int                   true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/status/block/{id} [put]
func (bh *BusinessHandler) HandleUpdateBlockBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: update block business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req struct {
		Status bool `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update block business status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = bh.businessStore.UpdateBlockBusinessStatus(id, req.Status)
	if err != nil {
		bh.logger.Printf("ERROR: update block business status: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business blocked status updated successfully"})
}

// HandleGetAllBusinesses godoc
// @Summary      Get all businesses
// @Description  Returns a list of all registered businesses
// @Tags         businesses
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of businesses"
// @Failure      404 {object} ErrorResponse "No businesses found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/get/all [get]
func (bh *BusinessHandler) HandleGetAllBusinesses(w http.ResponseWriter, r *http.Request) {
	res, err := bh.businessStore.GetAllBusinesses()
	if err != nil {
		bh.logger.Printf("ERROR: get all businesses: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if len(res) == 0 {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "no businesses found"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "businesses fetched successfully", "businesses": res})
}

// HandleDeleteBusiness godoc
// @Summary      Delete a business
// @Description  Deletes the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      404 {object} ErrorResponse "Business not found"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/delete/{id} [delete]
func (bh *BusinessHandler) HandleDeleteBusiness(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: delete business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = bh.businessStore.DeleteBusiness(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			bh.logger.Printf("ERROR: delete business: business not found: %v\n", err)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Printf("ERROR: delete business: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business deleted successfully"})
}

// HandleGetBusinessIDByUserID godoc
// @Summary      Get business ID by user ID
// @Description  Returns the business ID associated with the given user ID
// @Tags         businesses
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{} "business id"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      404 {object} ErrorResponse "No business found for this user"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/get/user/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessIDByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: get business id by user id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	businessID, err := bh.businessStore.GetBusinessIDByUserID(id)
	if err != nil {
		bh.logger.Printf("ERROR: get business id by user id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if businessID == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "no business found for this user"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business id fetched successfully", "business_id": businessID})
}

// HandleIsBusinessApproved godoc
// @Summary      Check if business is approved
// @Description  Returns whether the business with the given ID has been approved
// @Tags         businesses
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "approval status"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/status/{id} [get]
func (bh *BusinessHandler) HandleIsBusinessApproved(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		bh.logger.Printf("ERROR: is business approved: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	isApproved, err := bh.businessStore.IsBusinessApproved(id)
	if err != nil {
		bh.logger.Printf("ERROR: is business approved: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": isApproved})
}

func (bh *BusinessHandler) HandleRateBusiness(w http.ResponseWriter, r *http.Request) {
	var req businessRatingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: rate business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	err = bh.validateBusinessRatingRequest(&req)
	if err != nil {
		bh.logger.Printf("ERROR: rate business: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	ratings := &store.BusinessRating{
		BusinessID: req.BusinessID,
		UserID:     req.UserID,
		Rating:     req.Rating,
	}

	err = bh.businessStore.RateBusiness(ratings)
	if err != nil {
		bh.logger.Printf("ERROR: rate business: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
}
