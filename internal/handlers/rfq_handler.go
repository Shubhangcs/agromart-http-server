package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// rfqRequest represents the RFQ payload
type rfqRequest struct {
	ID            string    `json:"id"              example:"rfq-uuid-001"`
	BusinessID    string    `json:"business_id"     example:"biz-uuid-001"`
	CategoryID    string    `json:"category_id"     example:"cat-uuid-001"`
	SubCategoryID string    `json:"sub_category_id" example:"subcat-uuid-001"`
	ProductName   string    `json:"product_name"    example:"Wheat"`
	Quantity      float64   `json:"quantity"        example:"500"`
	Unit          string    `json:"unit"            example:"kg"`
	Price         float64   `json:"price"           example:"1200.50"`
	IsRFQActive   bool      `json:"is_rfq_active"   example:"true"`
	CreatedAT     time.Time `json:"created_at"`
	UpdatedAT     time.Time `json:"updated_at"`
}

// activateRFQRequest is used to toggle RFQ active status
type activateRFQRequest struct {
	IsRFQActive bool `json:"is_rfq_active" example:"true"`
}

type RFQHandler struct {
	rfqStore store.RFQStore
	logger   *log.Logger
}

func NewRFQHandler(rfqStore store.RFQStore, logger *log.Logger) *RFQHandler {
	return &RFQHandler{
		rfqStore: rfqStore,
		logger:   logger,
	}
}

func (rh *RFQHandler) validateCreateRFQ(req *rfqRequest) error {
	if req.BusinessID == "" {
		return errors.New("invalid request business id is required")
	}
	if req.CategoryID == "" {
		return errors.New("invalid request category id is required")
	}
	if req.SubCategoryID == "" {
		return errors.New("invalid request sub category id is required")
	}
	if req.ProductName == "" {
		return errors.New("invalid request product name is required")
	}
	if req.Quantity == 0 {
		return errors.New("invalid request product quantity is required")
	}
	if req.Unit == "" {
		return errors.New("invalid request unit is required")
	}
	if req.Price == 0 {
		return errors.New("invalid request product price is required")
	}
	return nil
}

// HandleCreateRFQ godoc
// @Summary      Create a new RFQ
// @Description  Creates a new Request for Quotation for a business
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        body body rfqRequest true "RFQ creation payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing required fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/create [post]
func (rh *RFQHandler) HandleCreateRFQ(w http.ResponseWriter, r *http.Request) {
	var req rfqRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rh.logger.Printf("ERROR: create rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = rh.validateCreateRFQ(&req)
	if err != nil {
		rh.logger.Printf("ERROR: create rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	rfq := &store.RFQ{
		BusinessID:    req.BusinessID,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		Price:         req.Price,
		IsRFQActive:   req.IsRFQActive,
	}

	err = rh.rfqStore.CreateRFQ(rfq)
	if err != nil {
		rh.logger.Printf("ERROR: create rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq created successfully"})
}

// HandleActivateRFQ godoc
// @Summary      Activate or deactivate an RFQ
// @Description  Toggles the active status of an RFQ by its ID
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        id   path int                true "RFQ ID"
// @Param        body body activateRFQRequest true "RFQ active status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/update/status/{id} [put]
func (rh *RFQHandler) HandleActivateRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		rh.logger.Printf("ERROR: activate rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	var req rfqRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rh.logger.Printf("ERROR: activate rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	rfq := &store.RFQ{
		ID:          id,
		IsRFQActive: req.IsRFQActive,
	}

	err = rh.rfqStore.ActivateRFQ(rfq)
	if err != nil {
		rh.logger.Printf("ERROR: activate rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq status changed successfully"})
}

// HandleUpdateRFQ godoc
// @Summary      Update an RFQ
// @Description  Updates the details of an existing RFQ by its ID
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        id   path int        true "RFQ ID"
// @Param        body body rfqRequest true "RFQ update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/update/{id} [put]
func (rh *RFQHandler) HandleUpdateRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		rh.logger.Printf("ERROR: update rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	var req rfqRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		rh.logger.Printf("ERROR: update rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	rfq := &store.RFQ{
		ID:            id,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		Price:         req.Price,
		Unit:          req.Unit,
	}

	err = rh.rfqStore.UpdateRFQ(rfq)
	if err != nil {
		rh.logger.Printf("ERROR: update rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq updated successfully"})
}

// HandleDeleteRFQ godoc
// @Summary      Delete an RFQ
// @Description  Deletes the RFQ with the given ID
// @Tags         rfq
// @Produce      json
// @Param        id path int true "RFQ ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/delete/{id} [delete]
func (rh *RFQHandler) HandleDeleteRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		rh.logger.Printf("ERROR: delete rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	err = rh.rfqStore.DeleteRFQ(id)
	if err != nil {
		rh.logger.Printf("ERROR: delete rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq deleted successfully"})
}

// HandleGetAllRFQ godoc
// @Summary      Get all RFQs
// @Description  Returns a list of all RFQs across all businesses
// @Tags         rfq
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of rfqs"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/get/all [get]
func (rh *RFQHandler) HandleGetAllRFQ(w http.ResponseWriter, r *http.Request) {
	res, err := rh.rfqStore.GetAllRFQ()
	if err != nil {
		rh.logger.Printf("ERROR: get all rfq: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfqs fetched successfully", "rfqs": res})
}

// HandleGetRFQByBusinessID godoc
// @Summary      Get RFQs by business ID
// @Description  Returns all RFQs belonging to the business with the given ID
// @Tags         rfq
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "list of rfqs"
// @Failure      400 {object} ErrorResponse "Invalid ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /rfq/get/{id} [get]
func (rh *RFQHandler) HandleGetRFQByBusinessID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		rh.logger.Printf("ERROR: get rfq by business id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := rh.rfqStore.GetRFQByBusinessID(id)
	if err != nil {
		rh.logger.Printf("ERROR: get rfq by business id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfqs fetched successfully", "rfqs": res})
}