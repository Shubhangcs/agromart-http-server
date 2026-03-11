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

type productRequest struct {
	ID              string                `json:"id"`
	BusinessID      string                `json:"business_id"`
	CategoryID      string                `json:"category_id"`
	SubCategoryID   string                `json:"sub_category_id"`
	Name            string                `json:"name"`
	Description     string                `json:"description"`
	Quantity        float64               `json:"quantity"`
	Unit            string                `json:"unit"`
	Price           float64               `json:"price"`
	MOQ             string                `json:"moq"`
	Images          []store.ProductImages `json:"product_images"`
	IsProductActive bool                  `json:"is_product_active"`
	CreatedAT       time.Time             `json:"created_at"`
	UpdatedAT       time.Time             `json:"updated_at"`
}

type ProductHandler struct {
	productStore store.ProductStore
	logger       *log.Logger
}

func NewProductHandler(productStore store.ProductStore, logger *log.Logger) *ProductHandler {
	return &ProductHandler{
		productStore: productStore,
		logger:       logger,
	}
}

func (ph *ProductHandler) validateCreateProductRequest(req *productRequest) error {
	if req.BusinessID == "" {
		return errors.New("invalid request business id is required")
	}
	if req.CategoryID == "" {
		return errors.New("invalid request category id is required")
	}
	if req.SubCategoryID == "" {
		return errors.New("invalid request sub category id is required")
	}
	if req.Name == "" {
		return errors.New("invalid request product name is required")
	}
	if req.Description == "" {
		return errors.New("invalid request product description is required")
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
	if req.MOQ == "" {
		return errors.New("invalid request moq is required")
	}
	return nil
}

func (ph *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var req productRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ph.logger.Printf("ERROR: create product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = ph.validateCreateProductRequest(&req)
	if err != nil {
		ph.logger.Printf("ERROR: create product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	product := &store.Product{
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

	err = ph.productStore.CreateProduct(product)
	if err != nil {
		ph.logger.Printf("ERROR: create product: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "product created successfully", "product_id": product.ID})
}

func (ph *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: update product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	var req productRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ph.logger.Printf("ERROR: update product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	product := &store.Product{
		ID:            id,
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		SubCategoryID: req.SubCategoryID,
		Quantity:      req.Quantity,
		Price:         req.Price,
		Unit:          req.Unit,
		MOQ:           req.MOQ,
	}

	err = ph.productStore.UpdateProduct(product)
	if err != nil {
		ph.logger.Printf("ERROR: update product: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product updated successfully"})
}

func (ph *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: delete product: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	err = ph.productStore.DeleteProduct(id)
	if err != nil {
		ph.logger.Printf("ERROR: delete product: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product deleted successfully"})
}

func (ph *ProductHandler) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	res, err := ph.productStore.GetAllProducts()
	if err != nil {
		ph.logger.Printf("ERROR: get all products: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "products fetched successfully", "products": res})
}

func (ph *ProductHandler) HandleGetBusinessProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: get business products: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := ph.productStore.GetBusinessProducts(id)
	if err != nil {
		ph.logger.Printf("ERROR: get business products: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "products fetched successfully", "products": res})
}

func (ph *ProductHandler) HandleGetCategoryBasedProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: get category based products: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := ph.productStore.GetCategoryBasedProducts(id)
	if err != nil {
		ph.logger.Printf("ERROR: get category based products: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "products fetched successfully", "products": res})
}

func (ph *ProductHandler) HandleGetSubCategoryBasedProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: get sub category based products: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := ph.productStore.GetSubCategoryBasedProducts(id)
	if err != nil {
		ph.logger.Printf("ERROR: get sub category based products: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "products fetched successfully", "products": res})
}

func (ph *ProductHandler) HandleGetFollowersProducts(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: get followers products: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := ph.productStore.GetFollowersProducts(id)
	if err != nil {
		ph.logger.Printf("ERROR: get followers products: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "products fetched successfully", "products": res})
}

func (ph *ProductHandler) HandleGetProductDetailsByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: get product details by id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	res, err := ph.productStore.GetProductDetailsByID(id)
	if err != nil {
		ph.logger.Printf("ERROR: get product details by id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product details fetched successfully", "product_details": res})
}

func (ph *ProductHandler) HandleChangeProductActivateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ph.logger.Printf("ERROR: change product active status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	var req productRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ph.logger.Printf("ERROR: change product active status: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	product := &store.Product{
		ID:              id,
		IsProductActive: req.IsProductActive,
	}

	err = ph.productStore.ChangeProductActivateStatus(product)
	if err != nil {
		ph.logger.Printf("ERROR: change product active status: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product status changed successfully"})
}
