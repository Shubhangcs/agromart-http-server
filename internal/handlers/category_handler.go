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

// categoryRequest represents the category payload
type categoryRequest struct {
	ID          string `json:"id"          example:"cat-uuid-001"`
	Name        string `json:"name"        example:"Grains"`
	Description string `json:"description" example:"All types of grains and cereals"`
	CreatedAT   string `json:"created_at"`
	UpdatedAT   string `json:"updated_at"`
}

// subCategoryRequest represents the sub-category payload
type subCategoryRequest struct {
	ID          string    `json:"id"           example:"subcat-uuid-001"`
	CategoryID  string    `json:"category_id"  example:"cat-uuid-001"`
	Name        string    `json:"name"         example:"Wheat"`
	Description string    `json:"description"  example:"All varieties of wheat"`
	CreatedAT   time.Time `json:"created_at"`
	UpdatedAT   time.Time `json:"updated_at"`
}

type CategoryHandler struct {
	categoryStore store.CategoryStore
	logger        *log.Logger
}

func NewCategoryHandler(categoryStore store.CategoryStore, logger *log.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryStore: categoryStore,
		logger:        logger,
	}
}

func (ch *CategoryHandler) validateCreateCategory(req *categoryRequest) error {
	if req.Name == "" {
		return errors.New("invalid request name is required")
	}
	if req.Description == "" {
		return errors.New("invalid request description is required")
	}
	return nil
}

func (ch *CategoryHandler) validateCreateSubCategory(req *subCategoryRequest) error {
	if req.CategoryID == "" {
		return errors.New("invalid request category id is required")
	}
	if req.Name == "" {
		return errors.New("invalid request name is required")
	}
	if req.Description == "" {
		return errors.New("invalid request description is required")
	}
	return nil
}

// HandleCreateCategory godoc
// @Summary      Create a category
// @Description  Creates a new product category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        body body categoryRequest true "Category creation payload"
// @Success      201 {object} map[string]interface{} "category created successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/create [post]
func (ch *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ch.logger.Printf("ERROR: create category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = ch.validateCreateCategory(&req)
	if err != nil {
		ch.logger.Printf("ERROR: create category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	category := &store.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	err = ch.categoryStore.CreateCategory(category)
	if err != nil {
		ch.logger.Printf("ERROR: create category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "category created successfully", "category_id": category.ID})
}

// HandleCreateSubCategory godoc
// @Summary      Create a sub-category
// @Description  Creates a new sub-category under an existing category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        body body subCategoryRequest true "Sub-category creation payload"
// @Success      201 {object} map[string]interface{} "sub category created successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/create [post]
func (ch *CategoryHandler) HandleCreateSubCategory(w http.ResponseWriter, r *http.Request) {
	var req subCategoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ch.logger.Printf("ERROR: create sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = ch.validateCreateSubCategory(&req)
	if err != nil {
		ch.logger.Printf("ERROR: create sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	subCategory := &store.SubCategory{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
	}
	err = ch.categoryStore.CreateSubCategory(subCategory)
	if err != nil {
		ch.logger.Printf("ERROR: create sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "sub category created successfully", "sub_category_id": subCategory.ID})
}

// HandleUpdateCategory godoc
// @Summary      Update a category
// @Description  Updates the name and description of the category with the given ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path int             true "Category ID"
// @Param        body body categoryRequest true "Category update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/update/{id} [put]
func (ch *CategoryHandler) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: update category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: update category: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	var req categoryRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ch.logger.Printf("ERROR: update category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	category := &store.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	err = ch.categoryStore.UpdateCategory(category)
	if err != nil {
		ch.logger.Printf("ERROR: update category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "category updated successfully"})
}

// HandleUpdateSubCategory godoc
// @Summary      Update a sub-category
// @Description  Updates the name and description of the sub-category with the given ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id   path int                true "Sub-category ID"
// @Param        body body subCategoryRequest true "Sub-category update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid ID or payload"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/update/{id} [put]
func (ch *CategoryHandler) HandleUpdateSubCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: update sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: update sub category: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	var req subCategoryRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ch.logger.Printf("ERROR: update sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	subCategory := &store.SubCategory{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	err = ch.categoryStore.UpdateSubCategory(subCategory)
	if err != nil {
		ch.logger.Printf("ERROR: update sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "sub category updated successfully"})
}

// HandleDeleteCategory godoc
// @Summary      Delete a category
// @Description  Deletes the category with the given ID
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid or missing ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/delete/{id} [delete]
func (ch *CategoryHandler) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: delete category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: delete category: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	err = ch.categoryStore.DeleteCategory(id)
	if err != nil {
		ch.logger.Printf("ERROR: delete category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "category deleted successfully"})
}

// HandleDeleteSubCategory godoc
// @Summary      Delete a sub-category
// @Description  Deletes the sub-category with the given ID
// @Tags         categories
// @Produce      json
// @Param        id path int true "Sub-category ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid or missing ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/delete/{id} [delete]
func (ch *CategoryHandler) HandleDeleteSubCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: delete sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: delete sub category: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	err = ch.categoryStore.DeleteSubCategory(id)
	if err != nil {
		ch.logger.Printf("ERROR: delete sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "sub category deleted successfully"})
}

// HandleGetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Returns the category with the given ID
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} map[string]interface{} "category details"
// @Failure      400 {object} ErrorResponse "Invalid or missing ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/get/{id} [get]
func (ch *CategoryHandler) HandleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: get category by id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: get category by id: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	res, err := ch.categoryStore.GetCategoryByID(id)
	if err != nil {
		ch.logger.Printf("ERROR: get category by id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "category fetched successfully", "category": res})
}

// HandleGetSubCategoryByID godoc
// @Summary      Get sub-category by ID
// @Description  Returns the sub-category with the given ID
// @Tags         categories
// @Produce      json
// @Param        id path int true "Sub-category ID"
// @Success      200 {object} map[string]interface{} "sub-category details"
// @Failure      400 {object} ErrorResponse "Invalid or missing ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/get/{id} [get]
func (ch *CategoryHandler) HandleGetSubCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: get sub category by id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: get sub category by id: invalid id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request id not found"})
		return
	}
	res, err := ch.categoryStore.GetSubCategoryByID(id)
	if err != nil {
		ch.logger.Printf("ERROR: get sub category by id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "sub category fetched successfully", "sub_category": res})
}

// HandleGetAllCategories godoc
// @Summary      Get all categories
// @Description  Returns a list of all product categories
// @Tags         categories
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of categories"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/get/all [get]
func (ch *CategoryHandler) HandleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	res, err := ch.categoryStore.GetAllCategories()
	if err != nil {
		ch.logger.Printf("ERROR: get all categories: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "categories fetched successfully", "categories": res})
}

// HandleGetAllSubCategories godoc
// @Summary      Get all sub-categories
// @Description  Returns a list of all sub-categories across all categories
// @Tags         categories
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of sub-categories"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/get/all [get]
func (ch *CategoryHandler) HandleGetAllSubCategories(w http.ResponseWriter, r *http.Request) {
	res, err := ch.categoryStore.GetAllSubCategories()
	if err != nil {
		ch.logger.Printf("ERROR: get all sub categories: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "sub categories fetched successfully", "sub_categories": res})
}

// HandleGetSubCategoriesByCategoryID godoc
// @Summary      Get sub-categories by category ID
// @Description  Returns all sub-categories belonging to the given category
// @Tags         categories
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} map[string]interface{} "list of sub-categories"
// @Failure      400 {object} ErrorResponse "Invalid or missing category ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/get/category/{id} [get]
func (ch *CategoryHandler) HandleGetSubCategoriesByCategoryID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		ch.logger.Printf("ERROR: get sub categories by category id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if id == "" {
		ch.logger.Println("ERROR: get sub categories by category id: invalid category id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request category id not found"})
		return
	}
	res, err := ch.categoryStore.GetSubCategoriesByCategoryID(id)
	if err != nil {
		ch.logger.Printf("ERROR: get sub categories by category id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "sub categories fetched successfully", "sub_categories": res})
}
