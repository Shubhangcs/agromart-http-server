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

// ProductHandler handles all product-related HTTP requests.
type ProductHandler struct {
	productStore store.ProductStore
	logger       *slog.Logger
}

func NewProductHandler(productStore store.ProductStore, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		productStore: productStore,
		logger:       logger,
	}
}

// HandleCreateProduct godoc
// @Summary      Create a product
// @Description  Creates a new product listing for a business
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body body models.CreateProductRequest true "Product creation payload"
// @Success      201 {object} map[string]interface{} "product_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/create [post]
func (ph *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ph.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	product := &models.Product{
		BusinessID:      req.BusinessID,
		CategoryID:      req.CategoryID,
		SubCategoryID:   req.SubCategoryID,
		Name:            req.Name,
		Description:     req.Description,
		Quantity:        req.Quantity,
		Unit:            req.Unit,
		Price:           req.Price,
		MOQ:             req.MOQ,
		IsProductActive: req.IsProductActive,
	}
	if err := ph.productStore.CreateProduct(product); err != nil {
		utils.ServerError(w, ph.logger, "create product", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "product created successfully", "product_id": product.ID})
}

// HandleUpdateProduct godoc
// @Summary      Update a product
// @Description  Updates the details of the product with the given ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path string                      true "Product ID"
// @Param        body body models.UpdateProductRequest true "Product update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/update/{id} [put]
func (ph *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	var req models.UpdateProductRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ph.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	if err = ph.productStore.UpdateProduct(&models.Product{
		ID:            id,
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		Quantity:      req.Quantity,
		Price:         req.Price,
		Unit:          req.Unit,
		MOQ:           req.MOQ,
		Description:   req.Description,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "product not found"})
			return
		}
		utils.ServerError(w, ph.logger, "update product", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product updated successfully"})
}

// HandleDeleteProduct godoc
// @Summary      Delete a product
// @Description  Deletes the product with the given ID
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/delete/{id} [delete]
func (ph *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	if err = ph.productStore.DeleteProduct(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "product not found"})
			return
		}
		utils.ServerError(w, ph.logger, "delete product", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product deleted successfully"})
}

// HandleGetAllProducts godoc
// @Summary      Get all products
// @Description  Returns a paginated list of active products. Supports optional name search and location filters.
// @Tags         products
// @Produce      json
// @Param        page  query string false "Page number (default 1)"
// @Param        limit query string false "Items per page (default 20, max 100)"
// @Param        q     query string false "Search by product name (case-insensitive)"
// @Param        city  query string false "Filter by business city (case-insensitive)"
// @Param        state query string false "Filter by business state (case-insensitive)"
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/all [get]
func (ph *ProductHandler) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	pg := utils.ReadPaginationParams(r)
	filter := utils.ReadProductFilter(r)
	res, err := ph.productStore.GetAllProducts(filter, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, ph.logger, "get all products", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "products fetched successfully",
		"products":   res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetBusinessProducts godoc
// @Summary      Get products by business
// @Description  Returns a paginated list of products for the given business ID
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "Business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/business/{id} [get]
func (ph *ProductHandler) HandleGetBusinessProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := ph.productStore.GetBusinessProducts(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, ph.logger, "get business products", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "products fetched successfully",
		"products":   res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetCategoryBasedProducts godoc
// @Summary      Get products by category
// @Description  Returns a paginated list of active products in the given category
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "Category ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/category/{id} [get]
func (ph *ProductHandler) HandleGetCategoryBasedProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := ph.productStore.GetCategoryBasedProducts(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, ph.logger, "get category based products", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "products fetched successfully",
		"products":   res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetSubCategoryBasedProducts godoc
// @Summary      Get products by sub-category
// @Description  Returns a paginated list of active products in the given sub-category
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "Sub-category ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/sub/category/{id} [get]
func (ph *ProductHandler) HandleGetSubCategoryBasedProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := ph.productStore.GetSubCategoryBasedProducts(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, ph.logger, "get sub category based products", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "products fetched successfully",
		"products":   res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetFollowersProducts godoc
// @Summary      Get products from followed businesses
// @Description  Returns a paginated list of active products from all businesses the user follows
// @Tags         products
// @Produce      json
// @Param        id    path  string true  "User ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/followers/{id} [get]
func (ph *ProductHandler) HandleGetFollowersProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := ph.productStore.GetFollowersProducts(id, pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, ph.logger, "get followers products", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "products fetched successfully",
		"products":   res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetProductDetailsByID godoc
// @Summary      Get product details
// @Description  Returns full product details including business, category, and image info
// @Tags         products
// @Produce      json
// @Param        id path string true "Product ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/get/{id} [get]
func (ph *ProductHandler) HandleGetProductDetailsByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	res, err := ph.productStore.GetProductDetailsByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "product not found"})
			return
		}
		utils.ServerError(w, ph.logger, "get product details by id", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product details fetched successfully", "product_details": res})
}

// HandleChangeProductActivateStatus godoc
// @Summary      Toggle product active status
// @Description  Activates or deactivates the product with the given ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path  string                            true "Product ID"
// @Param        body body  models.ChangeProductStatusRequest true "Status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /product/update/status/{id} [patch]
func (ph *ProductHandler) HandleChangeProductActivateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, ph.logger, err.Error(), err)
		return
	}
	var req models.ChangeProductStatusRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, ph.logger, "invalid request payload", err)
		return
	}
	if err = ph.productStore.ChangeProductActivateStatus(&models.Product{
		ID:              id,
		IsProductActive: req.IsProductActive,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "product not found"})
			return
		}
		utils.ServerError(w, ph.logger, "change product active status", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product status updated successfully"})
}
