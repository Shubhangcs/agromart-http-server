package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shubhangcs/agromart-server/internal/blob"
	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// PresignedURLResponse represents a successful presigned URL response
type PresignedURLResponse struct {
	URL     string `json:"url"     example:"https://s3.amazonaws.com/bucket/profile/admins/uuid.png?X-Amz-Signature=..."`
	Message string `json:"message" example:"presigned url generated successfully"`
}

type BlobHandler struct {
	blobStore store.BlobStore
	logger    *log.Logger
	blob      *blob.AWSS3
}

func NewBlobHandler(logger *log.Logger, blob *blob.AWSS3, blobStore store.BlobStore) *BlobHandler {
	return &BlobHandler{
		blobStore: blobStore,
		logger:    logger,
		blob:      blob,
	}
}

// HandleUpdateAdminProfileImage godoc
// @Summary      Generates presigned URL for admin profile image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload a profile image for the given admin
// @Tags         Image Upload
// @Produce      json
// @Param        id path uuid true "Admin ID"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid URL Param"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /admin/update/image/{id} [put]
func (bh *BlobHandler) HandleUpdateAdminProfileImage(w http.ResponseWriter, r *http.Request) {
	adminId, err := utils.ReadParamID(r)
	var time = time.Now().Unix()
	if err != nil {
		bh.logger.Printf("ERROR: update admin profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("profile/admins/%s_%d.png", adminId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update admin profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = bh.blobStore.UpdateAdminProfileImage(adminId, fmt.Sprintf("/profile/admins/%s_%d.png", adminId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update admin profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "admin profile presigned url generated successfully"})
}

// HandleUpdateUserProfileImage godoc
// @Summary      Generate presigned URL for user profile image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload a profile image for the given user
// @Tags         Image Upload
// @Produce      json
// @Param        id path uuid true "User ID"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid URL Param"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /user/update/image/{id} [put]
func (bh *BlobHandler) HandleUpdateUserProfileImage(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	var time = time.Now().Unix()
	if err != nil {
		bh.logger.Printf("ERROR: update user profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("profile/users/%s_%d.png", userId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update user profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = bh.blobStore.UpdateUserProfileImage(userId, fmt.Sprintf("/profile/users/%s_%d.png", userId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update user profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "user profile presigned url generated successfully"})
}

// HandleUpdateBusinessProfileImage godoc
// @Summary      Generate presigned URL for business profile image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload a profile image for the given business
// @Tags         Image Upload
// @Produce      json
// @Param        id path uuid true "Business ID"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid URL Param"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /business/update/image/{id} [put]
func (bh *BlobHandler) HandleUpdateBusinessProfileImage(w http.ResponseWriter, r *http.Request) {
	businessId, err := utils.ReadParamID(r)
	var time = time.Now().Unix()
	if err != nil {
		bh.logger.Printf("ERROR: update business profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("profile/business/%s_%d.png", businessId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update business profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = bh.blobStore.UpdateBusinessProfileImage(businessId, fmt.Sprintf("/profile/business/%s_%d.png", businessId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update business profile image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "business profile presigned url generated successfully"})
}

// HandleUpdateCategoryImage godoc
// @Summary      Generate presigned URL for category image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload an image for the given category
// @Tags         Image Upload
// @Produce      json
// @Param        id path uuid true "Category ID"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid URL Param"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/update/image/{id} [put]
func (bh *BlobHandler) HandleUpdateCategoryImage(w http.ResponseWriter, r *http.Request) {
	catId, err := utils.ReadParamID(r)
	var time = time.Now().Unix()
	if err != nil {
		bh.logger.Printf("ERROR: update category image: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("categories/%s_%d.png", catId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update category image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = bh.blobStore.UpdateCategoryImage(catId, fmt.Sprintf("/categories/%s_%d.png", catId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update category image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "category presigned url generated successfully"})
}

// HandleUpdateSubCategoryImage godoc
// @Summary      Generate presigned URL for sub-category image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload an image for the given sub-category
// @Tags         Image Upload
// @Produce      json
// @Param        id path uuid true "Sub-category ID"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid URL Param"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /category/sub/update/image/{id} [put]
func (bh *BlobHandler) HandleUpdateSubCategoryImage(w http.ResponseWriter, r *http.Request) {
	catId, err := utils.ReadParamID(r)
	var time = time.Now().Unix()
	if err != nil {
		bh.logger.Printf("ERROR: update sub category image: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("sub_categories/%s_%d.png", catId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update sub category image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = bh.blobStore.UpdateSubCategoryImage(catId, fmt.Sprintf("/sub_categories/%s_%d.png", catId, time))
	if err != nil {
		bh.logger.Printf("ERROR: update sub category image: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "sub category presigned url generated successfully"})
}

// productImageRequest represents the product image updation/deletion payload
type productImageRequest struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Index     int    `json:"index"`
	Image     string `json:"image"`
	CreatedAT int    `json:"created_at,omitempty"`
	UpdatedAT int    `json:"updated_at,omitempty"`
}

func (pir *productImageRequest) validateProductImageRequest() error {
	if pir.Index <= 0 || pir.Index > 3 {
		return errors.New("invalid request invalid index value")
	}

	if pir.ProductID == "" {
		return errors.New("invalid request product id is required")
	}

	return nil
}

func (pir *productImageRequest) validateDeleteProductImageRequest() error {
	if pir.ProductID == "" {
		return errors.New("invalid request product id is required")
	}

	if pir.ID == "" {
		return errors.New("invalid request image id is required")
	}

	return nil
}

// HandleUpdateProductImage godoc
// @Summary      Generate presigned URL for product image upload and updates the image path in database
// @Description  Returns an S3 presigned URL to upload an image for the given product
// @Tags         Image Upload
// @Accept       json
// @Produce      json
// @Param        body body productImageRequest true "Update Product Image Payload"
// @Success      200 {object} PresignedURLResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /product/update/image [put]
func (bh *BlobHandler) HandleUpdateProductImage(w http.ResponseWriter, r *http.Request) {
	var req productImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: update product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = req.validateProductImageRequest()
	if err != nil {
		bh.logger.Printf("ERROR: update product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	var id = uuid.NewString()
	productImage := &models.ProductImage{
		ID:        id,
		ProductID: req.ProductID,
		Index:     req.Index,
		Image:     fmt.Sprintf("/products/%s/%s.png", req.ProductID, id),
	}

	err = bh.blobStore.UpdateProductImage(productImage)
	if err != nil {
		bh.logger.Printf("ERROR: update product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	url, err := bh.blob.GenerateUploadPresignedURL(fmt.Sprintf("products/%s/%s.png", req.ProductID, id))
	if err != nil {
		bh.logger.Printf("ERROR: update product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product upload image presigned image generated successfully", "url": url})
}

// HandleDeleteProductImage godoc
// @Summary      Delete's the Product Image from S3 and image path from database
// @Description  Returns success message after successfully deleting image from both S3 and database
// @Tags         Image Upload
// @Accept       json
// @Produce      json
// @Param        body body productImageRequest true "Update Product Image Payload"
// @Success      200 {object} map[string]string{} "product image deleted successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /product/delete/image [delete]
func (bh *BlobHandler) HandleDeleteProductImage(w http.ResponseWriter, r *http.Request) {
	var req productImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		bh.logger.Printf("ERROR: delete product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}

	err = req.validateDeleteProductImageRequest()
	if err != nil {
		bh.logger.Printf("ERROR: delete product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	err = bh.blob.DeleteImage(fmt.Sprintf("products/%s/%s.png", req.ProductID, req.ID))
	if err != nil {
		bh.logger.Printf("ERROR: delete product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	err = bh.blobStore.DeleteProductImage(req.ID)
	if err != nil {
		bh.logger.Printf("ERROR: delete product image request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "product image deleted successfully"})
}
