package handlers

import (
	"encoding/json"
	"errors"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

// WishlistHandler handles all wishlist-related HTTP requests.
type WishlistHandler struct {
	wishlistStore store.WishlistStore
	logger        *slog.Logger
}

func NewWishlistHandler(wishlistStore store.WishlistStore, logger *slog.Logger) *WishlistHandler {
	return &WishlistHandler{wishlistStore: wishlistStore, logger: logger}
}

// HandleAddToWishlist godoc
// @Summary      Add product to wishlist
// @Description  Adds a product to the authenticated user's wishlist. Adding the same product twice is silently ignored.
// @Tags         wishlist
// @Accept       json
// @Produce      json
// @Param        body body models.AddToWishlistRequest true "Product to add"
// @Success      201 {object} handlers.MessageResponse
// @Failure      400 {object} handlers.ErrorResponse
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /wishlist/add [post]
func (wh *WishlistHandler) HandleAddToWishlist(w http.ResponseWriter, r *http.Request) {
	claims, ok := wh.claims(w, r)
	if !ok {
		return
	}

	var req models.AddToWishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, wh.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, wh.logger, err.Error(), err)
		return
	}

	if err := wh.wishlistStore.AddToWishlist(claims.UserID, req.ProductID); err != nil {
		utils.ServerError(w, wh.logger, "add to wishlist", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "product added to wishlist"})
}

// HandleRemoveFromWishlist godoc
// @Summary      Remove product from wishlist
// @Description  Removes a product from the authenticated user's wishlist.
// @Tags         wishlist
// @Produce      json
// @Param        product_id path string true "Product ID"
// @Success      200 {object} handlers.MessageResponse
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      404 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /wishlist/remove/{product_id} [delete]
func (wh *WishlistHandler) HandleRemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	claims, ok := wh.claims(w, r)
	if !ok {
		return
	}

	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		utils.BadRequest(w, wh.logger, "product_id is required", nil)
		return
	}

	if err := wh.wishlistStore.RemoveFromWishlist(claims.UserID, productID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "product not found in wishlist"})
			return
		}
		utils.ServerError(w, wh.logger, "remove from wishlist", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product removed from wishlist"})
}

// HandleGetWishlist godoc
// @Summary      Get user wishlist
// @Description  Returns all products saved in the authenticated user's wishlist, with product details.
// @Tags         wishlist
// @Produce      json
// @Param        page  query int false "Page number (default 1)"
// @Param        limit query int false "Items per page (default 20)"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /wishlist/get [get]
func (wh *WishlistHandler) HandleGetWishlist(w http.ResponseWriter, r *http.Request) {
	claims, ok := wh.claims(w, r)
	if !ok {
		return
	}

	pg := utils.ReadPaginationParams(r)
	items, err := wh.wishlistStore.GetUserWishlist(claims.UserID, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, wh.logger, "get wishlist", err)
		return
	}
	if items == nil {
		items = []models.WishlistItem{}
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"wishlist":   items,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleIsInWishlist godoc
// @Summary      Check if product is in wishlist
// @Description  Returns whether a specific product is already saved in the authenticated user's wishlist.
// @Tags         wishlist
// @Produce      json
// @Param        product_id path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} handlers.ErrorResponse
// @Failure      500 {object} handlers.ErrorResponse
// @Security     BearerAuth
// @Router       /wishlist/check/{product_id} [get]
func (wh *WishlistHandler) HandleIsInWishlist(w http.ResponseWriter, r *http.Request) {
	claims, ok := wh.claims(w, r)
	if !ok {
		return
	}

	productID := chi.URLParam(r, "product_id")
	if productID == "" {
		utils.BadRequest(w, wh.logger, "product_id is required", nil)
		return
	}

	exists, err := wh.wishlistStore.IsInWishlist(claims.UserID, productID)
	if err != nil {
		utils.ServerError(w, wh.logger, "check wishlist", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"in_wishlist": exists})
}

// claims extracts JWT claims from the request context, writing a 401 if missing.
func (wh *WishlistHandler) claims(w http.ResponseWriter, r *http.Request) (*tokens.Token, bool) {
	c, _ := r.Context().Value("claims").(*tokens.Token)
	if c == nil {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
		return nil, false
	}
	return c, true
}
