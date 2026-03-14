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

// ReviewHandler handles business and product review HTTP requests.
type ReviewHandler struct {
	reviewStore store.ReviewStore
	logger      *slog.Logger
}

func NewReviewHandler(reviewStore store.ReviewStore, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewStore: reviewStore,
		logger:      logger,
	}
}

// --- Business review handlers ---

// HandleCreateBusinessReview godoc
// @Summary      Create a business review
// @Description  Creates a new written review for a business
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        body body models.CreateBusinessReviewRequest true "Business review payload"
// @Success      201 {object} map[string]interface{} "review_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/review/create [post]
func (h *ReviewHandler) HandleCreateBusinessReview(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBusinessReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid request payload")
		return
	}
	if err := validator.Validate(&req); err != nil {
		badRequest(w, err.Error())
		return
	}
	rev := &models.BusinessReview{
		BusinessID: req.BusinessID,
		UserID:     req.UserID,
		Review:     req.Review,
	}
	if err := h.reviewStore.CreateBusinessReview(rev); err != nil {
		serverError(w, h.logger, "create business review", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"message":   "business review created successfully",
		"review_id": rev.ID,
	})
}

// HandleUpdateBusinessReview godoc
// @Summary      Update a business review
// @Description  Updates the review text of the business review with the given ID
// @Tags         businesses
// @Accept       json
// @Produce      json
// @Param        id   path string                     true "Review ID"
// @Param        body body models.UpdateReviewRequest true "Review update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/review/update/{id} [put]
func (h *ReviewHandler) HandleUpdateBusinessReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	var req models.UpdateReviewRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid request payload")
		return
	}
	if err = validator.Validate(&req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if err = h.reviewStore.UpdateBusinessReview(&models.BusinessReview{ID: id, Review: req.Review}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "review not found"})
			return
		}
		serverError(w, h.logger, "update business review", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business review updated successfully"})
}

// HandleDeleteBusinessReview godoc
// @Summary      Delete a business review
// @Description  Deletes the business review with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id path string true "Review ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/review/delete/{id} [delete]
func (h *ReviewHandler) HandleDeleteBusinessReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	if err = h.reviewStore.DeleteBusinessReview(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "review not found"})
			return
		}
		serverError(w, h.logger, "delete business review", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "business review deleted successfully"})
}

// HandleGetBusinessReviews godoc
// @Summary      Get all business reviews
// @Description  Returns a paginated list of reviews for the business with the given ID
// @Tags         businesses
// @Produce      json
// @Param        id    path  string true  "Business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /business/review/get/{id} [get]
func (h *ReviewHandler) HandleGetBusinessReviews(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := h.reviewStore.GetBusinessReviews(id, pg.Limit, pg.Offset())
	if err != nil {
		serverError(w, h.logger, "get business reviews", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "business reviews fetched successfully",
		"reviews":    res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// --- Product review handlers ---

// HandleCreateProductReview godoc
// @Summary      Create a product review
// @Description  Creates a new written review for a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body body models.CreateProductReviewRequest true "Product review payload"
// @Success      201 {object} map[string]interface{} "review_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/review/create [post]
func (h *ReviewHandler) HandleCreateProductReview(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid request payload")
		return
	}
	if err := validator.Validate(&req); err != nil {
		badRequest(w, err.Error())
		return
	}
	rev := &models.ProductReview{
		ProductID: req.ProductID,
		UserID:    req.UserID,
		Review:    req.Review,
	}
	if err := h.reviewStore.CreateProductReview(rev); err != nil {
		serverError(w, h.logger, "create product review", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{
		"message":   "product review created successfully",
		"review_id": rev.ID,
	})
}

// HandleUpdateProductReview godoc
// @Summary      Update a product review
// @Description  Updates the review text of the product review with the given ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path string                     true "Review ID"
// @Param        body body models.UpdateReviewRequest true "Review update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/review/update/{id} [put]
func (h *ReviewHandler) HandleUpdateProductReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	var req models.UpdateReviewRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		badRequest(w, "invalid request payload")
		return
	}
	if err = validator.Validate(&req); err != nil {
		badRequest(w, err.Error())
		return
	}
	if err = h.reviewStore.UpdateProductReview(&models.ProductReview{ID: id, Review: req.Review}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "review not found"})
			return
		}
		serverError(w, h.logger, "update product review", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product review updated successfully"})
}

// HandleDeleteProductReview godoc
// @Summary      Delete a product review
// @Description  Deletes the product review with the given ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Review ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/review/delete/{id} [delete]
func (h *ReviewHandler) HandleDeleteProductReview(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	if err = h.reviewStore.DeleteProductReview(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "review not found"})
			return
		}
		serverError(w, h.logger, "delete product review", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product review deleted successfully"})
}

// HandleGetProductReviews godoc
// @Summary      Get all product reviews
// @Description  Returns a paginated list of reviews for the product with the given ID
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "Product ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/review/get/{id} [get]
func (h *ReviewHandler) HandleGetProductReviews(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := h.reviewStore.GetProductReviews(id, pg.Limit, pg.Offset())
	if err != nil {
		serverError(w, h.logger, "get product reviews", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "product reviews fetched successfully",
		"reviews":    res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}
