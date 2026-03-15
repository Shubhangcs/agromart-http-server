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

// RatingHandler handles product and business rating HTTP requests.
type RatingHandler struct {
	ratingStore store.RatingStore
	logger      *slog.Logger
}

func NewRatingHandler(ratingStore store.RatingStore, logger *slog.Logger) *RatingHandler {
	return &RatingHandler{
		ratingStore: ratingStore,
		logger:      logger,
	}
}

// HandleRateProduct godoc
// @Summary      Rate a product
// @Description  Submits or updates a rating for a product by a user (upsert). Rating must be between 0.5 and 5.0.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body body models.RateProductRequest true "Product rating payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/rate [post]
func (h *RatingHandler) HandleRateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.RateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	rating := &models.ProductRating{
		ProductID: req.ProductID,
		UserID:    req.UserID,
		Rating:    req.Rating,
	}
	if err := h.ratingStore.RateProduct(rating); err != nil {
		utils.ServerError(w, h.logger, "rate product", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":   "product rated successfully",
		"rating_id": rating.ID,
	})
}

// HandleGetAverageProductRating godoc
// @Summary      Get average product rating
// @Description  Returns the average rating for the product with the given ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/rate/average/{id} [get]
func (h *RatingHandler) HandleGetAverageProductRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	avg, err := h.ratingStore.GetAverageProductRating(id)
	if err != nil {
		utils.ServerError(w, h.logger, "get average product rating", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":        "average product rating fetched successfully",
		"average_rating": avg,
	})
}

// HandleGetProductRatings godoc
// @Summary      Get all product ratings
// @Description  Returns a paginated list of ratings for the product with the given ID
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "Product ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/rate/get/{id} [get]
func (h *RatingHandler) HandleGetProductRatings(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := h.ratingStore.GetRatingsByProductID(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, h.logger, "get product ratings", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "product ratings fetched successfully",
		"ratings":    res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleDeleteProductRating godoc
// @Summary      Delete a product rating
// @Description  Deletes the product rating with the given ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Rating ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/rate/delete/{id} [delete]
func (h *RatingHandler) HandleDeleteProductRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	if err = h.ratingStore.DeleteProductRating(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "rating not found"})
			return
		}
		utils.ServerError(w, h.logger, "delete product rating", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product rating deleted successfully"})
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
func (h *RatingHandler) HandleRateBusiness(w http.ResponseWriter, r *http.Request) {
	var req models.RateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	if err := h.ratingStore.RateBusiness(&models.BusinessRating{
		BusinessID: req.BusinessID,
		UserID:     req.UserID,
		Rating:     req.Rating,
	}); err != nil {
		utils.ServerError(w, h.logger, "rate business", err)
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
func (h *RatingHandler) HandleGetAverageBusinessRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	avg, err := h.ratingStore.GetAverageBusinessRating(id)
	if err != nil {
		utils.ServerError(w, h.logger, "get average business rating", err)
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
func (h *RatingHandler) HandleGetBusinessRatings(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := h.ratingStore.GetRatingsByBusinessID(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, h.logger, "get business ratings", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "business ratings fetched successfully",
		"ratings":    res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}
