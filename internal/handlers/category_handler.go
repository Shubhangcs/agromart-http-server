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

// CategoryHandler handles all category and sub-category HTTP requests.
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

func (ch *CategoryHandler) validateCreateCategory(req *models.CreateCategoryRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	return nil
}

func (ch *CategoryHandler) validateCreateSubCategory(req *models.CreateSubCategoryRequest) error {
	if req.CategoryID == "" {
		return errors.New("category_id is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	return nil
}

// HandleCreateCategory godoc
// @Summary      Create a category
// @Description  Creates a new product category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        body body models.CreateCategoryRequest true "Category creation payload"
// @Success      201 {object} map[string]interface{} "category_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/create [post]
func (ch *CategoryHandler) HandleCreateCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.Printf("ERROR: create category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := ch.validateCreateCategory(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := ch.categoryStore.CreateCategory(category); err != nil {
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
// @Param        body body models.CreateSubCategoryRequest true "Sub-category creation payload"
// @Success      201 {object} map[string]interface{} "sub_category_id returned"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/sub/create [post]
func (ch *CategoryHandler) HandleCreateSubCategory(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSubCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.Printf("ERROR: create sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := ch.validateCreateSubCategory(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	subCategory := &models.SubCategory{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := ch.categoryStore.CreateSubCategory(subCategory); err != nil {
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
// @Param        id   path  string                       true "Category ID"
// @Param        body body  models.UpdateCategoryRequest true "Category update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/update/{id} [put]
func (ch *CategoryHandler) HandleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateCategoryRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.Printf("ERROR: update category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = ch.categoryStore.UpdateCategory(&models.Category{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "category not found"})
			return
		}
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
// @Param        id   path  string                          true "Sub-category ID"
// @Param        body body  models.UpdateSubCategoryRequest true "Sub-category update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/sub/update/{id} [put]
func (ch *CategoryHandler) HandleUpdateSubCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateSubCategoryRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		ch.logger.Printf("ERROR: update sub category: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = ch.categoryStore.UpdateSubCategory(&models.SubCategory{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "sub category not found"})
			return
		}
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
// @Param        id path string true "Category ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/delete/{id} [delete]
func (ch *CategoryHandler) HandleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = ch.categoryStore.DeleteCategory(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "category not found"})
			return
		}
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
// @Param        id path string true "Sub-category ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/sub/delete/{id} [delete]
func (ch *CategoryHandler) HandleDeleteSubCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = ch.categoryStore.DeleteSubCategory(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "sub category not found"})
			return
		}
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
// @Param        id path string true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/get/{id} [get]
func (ch *CategoryHandler) HandleGetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := ch.categoryStore.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "category not found"})
			return
		}
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
// @Param        id path string true "Sub-category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/sub/get/{id} [get]
func (ch *CategoryHandler) HandleGetSubCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := ch.categoryStore.GetSubCategoryByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "sub category not found"})
			return
		}
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
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} ErrorResponse
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
// @Success      200 {object} map[string]interface{}
// @Failure      500 {object} ErrorResponse
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
// @Param        id path string true "Category ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /category/sub/get/category/{id} [get]
func (ch *CategoryHandler) HandleGetSubCategoriesByCategoryID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
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
