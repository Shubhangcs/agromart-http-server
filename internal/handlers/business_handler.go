package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

// BusinessHandler handles all business-related HTTP requests.
type BusinessHandler struct {
	businessStore store.BusinessStore
	logger        *slog.Logger
}

func NewBusinessHandler(businessStore store.BusinessStore, logger *slog.Logger) *BusinessHandler {
	return &BusinessHandler{
		businessStore: businessStore,
		logger:        logger,
	}
}

func (bh *BusinessHandler) HandleCreateBusiness(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("create business", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	business := &models.Business{
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
	if err := bh.businessStore.CreateBusiness(business); err != nil {
		bh.logger.Error("create business", "error", err)
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
// @Param        body body models.CreateSocialRequest true "Social links payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/social/create [post]
func (bh *BusinessHandler) HandleCreateSocial(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSocialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("create social", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	social := &models.Social{
		ID:        req.ID,
		Linkedin:  req.Linkedin,
		Instagram: req.Instagram,
		Facebook:  req.Facebook,
		Youtube:   req.Youtube,
		Telegram:  req.Telegram,
		X:         req.X,
		Website:   req.Website,
	}
	if err := bh.businessStore.CreateSocial(social); err != nil {
		bh.logger.Error("create social", "error", err)
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
// @Param        body body models.CreateLegalRequest true "Legal documents payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/legal/create [post]
func (bh *BusinessHandler) HandleCreateLegal(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLegalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("create legal", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	legal := &models.Legal{
		ID:           req.ID,
		Aadhaar:      req.Aadhaar,
		Pan:          req.Pan,
		ExportImport: req.ExportImport,
		MSME:         req.MSME,
		Fassi:        req.Fassi,
		GST:          req.GST,
	}
	if err := bh.businessStore.CreateLegal(legal); err != nil {
		bh.logger.Error("create legal", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "business legals created successfully"})
}

// HandleCreateBusinessApplication godoc
// @Summary      Submit a business application
// @Description  Submits an application for business approval
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body models.CreateApplicationRequest true "Business application payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/application/create [post]
func (bh *BusinessHandler) HandleCreateBusinessApplication(w http.ResponseWriter, r *http.Request) {
	var req models.CreateApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("create business application", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err := bh.businessStore.CreateBusinessApplication(&models.BusinessApplication{ID: req.ID, Status: "APPLIED"}); err != nil {
		bh.logger.Error("create business application", "error", err)
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
// @Param        id   path string                       true "Business ID"
// @Param        body body models.UpdateBusinessRequest true "Business update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateBusiness(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateBusinessRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("update business", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.UpdateBusiness(&models.Business{
		ID:           id,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		City:         req.City,
		State:        req.State,
		Pincode:      req.Pincode,
		BusinessType: req.BusinessType,
	}); err != nil {
		bh.logger.Error("update business", "error", err)
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
// @Param        id   path string                      true "Business ID"
// @Param        body body models.CreateSocialRequest  true "Social links update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/social/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateSocials(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.CreateSocialRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("update socials", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.UpdateSocial(&models.Social{
		ID:        id,
		Linkedin:  req.Linkedin,
		Instagram: req.Instagram,
		Youtube:   req.Youtube,
		Telegram:  req.Telegram,
		X:         req.X,
		Facebook:  req.Facebook,
		Website:   req.Website,
	}); err != nil {
		bh.logger.Error("update socials", "error", err)
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
// @Param        id   path string                     true "Business ID"
// @Param        body body models.CreateLegalRequest  true "Legal documents update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/legal/update/{id} [put]
func (bh *BusinessHandler) HandleUpdateLegals(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.CreateLegalRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("update legals", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = bh.businessStore.UpdateLegal(&models.Legal{
		ID:           id,
		Aadhaar:      req.Aadhaar,
		Pan:          req.Pan,
		ExportImport: req.ExportImport,
		MSME:         req.MSME,
		Fassi:        req.Fassi,
		GST:          req.GST,
	}); err != nil {
		bh.logger.Error("update legals", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/application/accept/{id} [put]
func (bh *BusinessHandler) HandleAcceptBusinessApplication(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.AcceptBusinessApplication(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Error("accept business application", "error", err)
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
// @Param        id   path string                           true "Business ID"
// @Param        body body models.RejectApplicationRequest  true "Rejection reason payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/application/reject/{id} [put]
func (bh *BusinessHandler) HandleRejectBusinessApplication(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.RejectApplicationRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("reject business application", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.RejectBusinessApplication(&models.BusinessApplication{
		ID:           id,
		Status:       "REJECTED",
		RejectReason: &req.RejectReason,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Error("reject business application", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/get/complete/{id} [get]
func (bh *BusinessHandler) HandleGetCompleteBusinessDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetCompleteBusinessDetails(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Error("get complete business details", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/get/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetBusiness(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Error("get business details", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/social/get/{id} [get]
func (bh *BusinessHandler) HandleGetSocialDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetSocial(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "social details not found"})
			return
		}
		bh.logger.Error("get business social details", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/legal/get/{id} [get]
func (bh *BusinessHandler) HandleGetLegalDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetLegal(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "legal details not found"})
			return
		}
		bh.logger.Error("get business legal details", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/application/get/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessApplicationDetails(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetBusinessApplication(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "application not found"})
			return
		}
		bh.logger.Error("get business application details", "error", err)
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
// @Param        id   path string                           true "Business ID"
// @Param        body body models.UpdateBusinessStatusRequest true "Status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/status/verify/{id} [put]
func (bh *BusinessHandler) HandleUpdateVerifyBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateBusinessStatusRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.UpdateVerifyBusinessStatus(id, req.Status); err != nil {
		bh.logger.Error("update verify business status", "error", err)
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
// @Param        id   path string                           true "Business ID"
// @Param        body body models.UpdateBusinessStatusRequest true "Status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/status/trust/{id} [put]
func (bh *BusinessHandler) HandleUpdateTrustBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateBusinessStatusRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = bh.businessStore.UpdateTrustBusinessStatus(id, req.Status); err != nil {
		bh.logger.Error("update trust business status", "error", err)
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
// @Param        id   path string                           true "Business ID"
// @Param        body body models.UpdateBusinessStatusRequest true "Status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/status/block/{id} [put]
func (bh *BusinessHandler) HandleUpdateBlockBusinessStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateBusinessStatusRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = bh.businessStore.UpdateBlockBusinessStatus(id, req.Status); err != nil {
		bh.logger.Error("update block business status", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business blocked status updated successfully"})
}

// HandleGetAllBusinesses godoc
// @Summary      Get all businesses
// @Description  Returns a paginated list of all registered businesses
// @Tags         businesses
// @Produce      json
// @Param        page  query int false "Page number (default 1)"
// @Param        limit query int false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/get/all [get]
func (bh *BusinessHandler) HandleGetAllBusinesses(w http.ResponseWriter, r *http.Request) {
	pg := utils.ReadPaginationParams(r)
	res, err := bh.businessStore.GetAllBusinesses(pg.Limit, pg.Offset())
	if err != nil {
		bh.logger.Error("get all businesses", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "businesses fetched successfully",
		"businesses": res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleDeleteBusiness godoc
// @Summary      Delete a business
// @Description  Deletes the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path string true "Business ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/delete/{id} [delete]
func (bh *BusinessHandler) HandleDeleteBusiness(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = bh.businessStore.DeleteBusiness(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "business not found"})
			return
		}
		bh.logger.Error("delete business", "error", err)
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
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/get/user/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessIDByUserID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	businessID, err := bh.businessStore.GetBusinessIDByUserID(id)
	if err != nil {
		bh.logger.Error("get business id by user id", "error", err)
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
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/status/{id} [get]
func (bh *BusinessHandler) HandleIsBusinessApproved(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	isApproved, err := bh.businessStore.IsBusinessApproved(id)
	if err != nil {
		bh.logger.Error("is business approved", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": isApproved})
}

// HandleRateBusiness godoc
// @Summary      Rate a business
// @Description  Submits or updates a rating for a business by a user (upsert)
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body models.RateBusinessRequest true "Rating payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/rate [post]
func (bh *BusinessHandler) HandleRateBusiness(w http.ResponseWriter, r *http.Request) {
	var req models.RateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		bh.logger.Error("rate business", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err := bh.businessStore.RateBusiness(&models.BusinessRating{
		BusinessID: req.BusinessID,
		UserID:     req.UserID,
		Rating:     req.Rating,
	}); err != nil {
		bh.logger.Error("rate business", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business rated successfully"})
}

// HandleGetAverageBusinessRating godoc
// @Summary      Get average business rating
// @Description  Returns the average rating for the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/rate/average/{id} [get]
func (bh *BusinessHandler) HandleGetAverageBusinessRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	avg, err := bh.businessStore.GetAverageBusinessRating(id)
	if err != nil {
		bh.logger.Error("get average business rating", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":        "average business rating fetched successfully",
		"average_rating": avg,
	})
}

// HandleGetBusinessRatings godoc
// @Summary      Get all business ratings
// @Description  Returns a list of all user ratings for the given business
// @Tags         businesses
// @Produce      json
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/rate/get/{id} [get]
func (bh *BusinessHandler) HandleGetBusinessRatings(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := bh.businessStore.GetRatingsByBusinessID(id)
	if err != nil {
		bh.logger.Error("get business ratings", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message": "business ratings fetched successfully",
		"ratings": res,
	})
}
