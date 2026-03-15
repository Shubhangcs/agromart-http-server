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

// RFQHandler handles all RFQ-related HTTP requests.
type RFQHandler struct {
	rfqStore store.RFQStore
	logger   *slog.Logger
}

func NewRFQHandler(rfqStore store.RFQStore, logger *slog.Logger) *RFQHandler {
	return &RFQHandler{
		rfqStore: rfqStore,
		logger:   logger,
	}
}

// HandleCreateRFQ godoc
// @Summary      Create a new RFQ
// @Description  Creates a new Request for Quotation for a business
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        body body models.CreateRFQRequest true "RFQ creation payload"
// @Success      201 {object} map[string]interface{} "rfq_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/create [post]
func (rh *RFQHandler) HandleCreateRFQ(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRFQRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, rh.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	rfq := &models.RFQ{
		BusinessID:    req.BusinessID,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		Unit:          req.Unit,
		Price:         req.Price,
		IsRFQActive:   req.IsRFQActive,
	}
	if err := rh.rfqStore.CreateRFQ(rfq); err != nil {
		utils.ServerError(w, rh.logger, "create rfq", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "rfq created successfully", "rfq_id": rfq.ID})
}

// HandleActivateRFQ godoc
// @Summary      Activate or deactivate an RFQ
// @Description  Toggles the active status of an RFQ by its ID
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        id   path  string                    true "RFQ ID"
// @Param        body body  models.ActivateRFQRequest true "RFQ active status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/update/status/{id} [put]
func (rh *RFQHandler) HandleActivateRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	var req models.ActivateRFQRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, rh.logger, "invalid request payload", err)
		return
	}
	if err = rh.rfqStore.ActivateRFQ(&models.RFQ{ID: id, IsRFQActive: req.IsRFQActive}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "rfq not found"})
			return
		}
		utils.ServerError(w, rh.logger, "activate rfq", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq status updated successfully"})
}

// HandleUpdateRFQ godoc
// @Summary      Update an RFQ
// @Description  Updates the details of an existing RFQ by its ID
// @Tags         rfq
// @Accept       json
// @Produce      json
// @Param        id   path  string                  true "RFQ ID"
// @Param        body body  models.UpdateRFQRequest true "RFQ update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/update/{id} [put]
func (rh *RFQHandler) HandleUpdateRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	var req models.UpdateRFQRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, rh.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	if err = rh.rfqStore.UpdateRFQ(&models.RFQ{
		ID:            id,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		ProductName:   req.ProductName,
		Quantity:      req.Quantity,
		Price:         req.Price,
		Unit:          req.Unit,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "rfq not found"})
			return
		}
		utils.ServerError(w, rh.logger, "update rfq", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq updated successfully"})
}

// HandleDeleteRFQ godoc
// @Summary      Delete an RFQ
// @Description  Deletes the RFQ with the given ID
// @Tags         rfq
// @Produce      json
// @Param        id path string true "RFQ ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/delete/{id} [delete]
func (rh *RFQHandler) HandleDeleteRFQ(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	if err = rh.rfqStore.DeleteRFQ(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "rfq not found"})
			return
		}
		utils.ServerError(w, rh.logger, "delete rfq", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "rfq deleted successfully"})
}

// HandleGetAllRFQ godoc
// @Summary      Get all RFQs
// @Description  Returns a paginated list of all RFQs across all businesses. Supports optional search and location filters.
// @Tags         rfq
// @Produce      json
// @Param        page  query string false "Page number (default 1)"
// @Param        limit query string false "Items per page (default 20, max 100)"
// @Param        q     query string false "Search by product name (case-insensitive)"
// @Param        city  query string false "Filter by business city (case-insensitive)"
// @Param        state query string false "Filter by business state (case-insensitive)"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/get/all [get]
func (rh *RFQHandler) HandleGetAllRFQ(w http.ResponseWriter, r *http.Request) {
	pg := utils.ReadPaginationParams(r)
	filter := utils.ReadRFQFilter(r)
	res, err := rh.rfqStore.GetAllRFQ(filter, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, rh.logger, "get all rfq", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "rfqs fetched successfully",
		"rfqs":       res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetRFQByBusinessID godoc
// @Summary      Get RFQs by business ID
// @Description  Returns a paginated list of RFQs belonging to the given business
// @Tags         rfq
// @Produce      json
// @Param        id    path  string true  "Business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /rfq/get/{id} [get]
func (rh *RFQHandler) HandleGetRFQByBusinessID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, rh.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := rh.rfqStore.GetRFQByBusinessID(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, rh.logger, "get rfq by business id", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "rfqs fetched successfully",
		"rfqs":       res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}
