package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

type LeadHandler struct {
	leadStore store.LeadStore
	logger    *slog.Logger
}

func NewLeadHandler(leadStore store.LeadStore, logger *slog.Logger) *LeadHandler {
	return &LeadHandler{leadStore: leadStore, logger: logger}
}

// HandleCreateLead godoc
// @Summary      Submit a product enquiry
// @Description  Creates a lead (enquiry) from one business to another for a specific product.
// @Tags         leads
// @Accept       json
// @Produce      json
// @Param        body body models.CreateLeadRequest true "Lead enquiry payload"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} handlers.ErrorResponse
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /leads/create [post]
func (lh *LeadHandler) HandleCreateLead(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, lh.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, lh.logger, err.Error(), err)
		return
	}
	if req.EnquirerBusinessID == req.EnquireToID {
		utils.BadRequest(w, lh.logger, "cannot enquire about your own product", nil)
		return
	}

	lead := &models.Lead{
		EnquirerBusinessID: req.EnquirerBusinessID,
		EnquireToID:        req.EnquireToID,
		ProductID:          req.ProductID,
		EnquiryMessage:     req.EnquiryMessage,
		ExpectedPrice:      req.ExpectedPrice,
	}

	if err := lh.leadStore.CreateLead(lead); err != nil {
		utils.ServerError(w, lh.logger, "create lead", err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"message": "enquiry submitted successfully",
		"lead_id": lead.LeadID,
	})
}

// HandleGetSentLeads godoc
// @Summary      Get sent enquiries
// @Description  Returns all leads submitted by the given business, with product and seller business details.
// @Tags         leads
// @Produce      json
// @Param        id    path  string true  "Enquirer business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20)"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /leads/sent/{id} [get]
func (lh *LeadHandler) HandleGetSentLeads(w http.ResponseWriter, r *http.Request) {
	businessID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, lh.logger, err.Error(), err)
		return
	}

	pg := utils.ReadPaginationParams(r)
	leads, err := lh.leadStore.GetSentLeads(businessID, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, lh.logger, "get sent leads", err)
		return
	}
	if leads == nil {
		leads = []models.LeadSentResponse{}
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"leads":      leads,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetReceivedLeads godoc
// @Summary      Get received enquiries
// @Description  Returns all leads received by the given business, with product and enquirer business details.
// @Tags         leads
// @Produce      json
// @Param        id    path  string true  "Seller business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20)"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /leads/received/{id} [get]
func (lh *LeadHandler) HandleGetReceivedLeads(w http.ResponseWriter, r *http.Request) {
	businessID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, lh.logger, err.Error(), err)
		return
	}

	pg := utils.ReadPaginationParams(r)
	leads, err := lh.leadStore.GetReceivedLeads(businessID, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, lh.logger, "get received leads", err)
		return
	}
	if leads == nil {
		leads = []models.LeadReceivedResponse{}
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"leads":      leads,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}
