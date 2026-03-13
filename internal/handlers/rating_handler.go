package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// ProductRatingHandler handles product rating HTTP requests.
type ProductRatingHandler struct {
	ratingStore store.ProductRatingStore
	logger      *log.Logger
}

func NewProductRatingHandler(ratingStore store.ProductRatingStore, logger *log.Logger) *ProductRatingHandler {
	return &ProductRatingHandler{
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
func (h *ProductRatingHandler) HandleRateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.RateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: rate product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if req.ProductID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "product_id is required"})
		return
	}
	if req.UserID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "user_id is required"})
		return
	}
	if req.Rating < 0.5 || req.Rating > 5.0 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "rating must be between 0.5 and 5.0"})
		return
	}
	rating := &models.ProductRating{
		ProductID: req.ProductID,
		UserID:    req.UserID,
		Rating:    req.Rating,
	}
	if err := h.ratingStore.RateProduct(rating); err != nil {
		h.logger.Printf("ERROR: rate product: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
func (h *ProductRatingHandler) HandleGetAverageProductRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	avg, err := h.ratingStore.GetAverageProductRating(id)
	if err != nil {
		h.logger.Printf("ERROR: get average product rating: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
func (h *ProductRatingHandler) HandleGetProductRatings(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := h.ratingStore.GetRatingsByProductID(id, pg.Limit, pg.Offset())
	if err != nil {
		h.logger.Printf("ERROR: get product ratings: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
func (h *ProductRatingHandler) HandleDeleteProductRating(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = h.ratingStore.DeleteProductRating(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "rating not found"})
			return
		}
		h.logger.Printf("ERROR: delete product rating: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product rating deleted successfully"})
}
