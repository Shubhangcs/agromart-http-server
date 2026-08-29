package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/shubhangcs/agromart-server/internal/blob"
	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

type BannerHandler struct {
	bannerStore store.BannerStore
	blob        *blob.AWSS3
	logger      *slog.Logger
}

func NewBannerHandler(bannerStore store.BannerStore, b *blob.AWSS3, logger *slog.Logger) *BannerHandler {
	return &BannerHandler{bannerStore: bannerStore, blob: b, logger: logger}
}

// HandleGetActiveBanners godoc
// @Summary      Home-screen banners (active, with image, in display order)
// @Tags         banners
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /banners/get/active [get]
func (h *BannerHandler) HandleGetActiveBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerStore.GetActive()
	if err != nil {
		utils.ServerError(w, h.logger, "get banners", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "banners fetched successfully", "banners": banners})
}

// HandleGetAllBanners godoc
// @Summary      All banners including inactive (admin)
// @Tags         banners
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /banners/get/all [get]
func (h *BannerHandler) HandleGetAllBanners(w http.ResponseWriter, r *http.Request) {
	banners, err := h.bannerStore.GetAll()
	if err != nil {
		utils.ServerError(w, h.logger, "get all banners", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "banners fetched successfully", "banners": banners})
}

// HandleCreateBanner godoc
// @Summary      Create a banner (admin); upload the image with /banners/update/image/{id}
// @Tags         banners
// @Accept       json
// @Produce      json
// @Param        body body models.UpsertBannerRequest true "Banner"
// @Success      201 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /banners/create [post]
func (h *BannerHandler) HandleCreateBanner(w http.ResponseWriter, r *http.Request) {
	var req models.UpsertBannerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	b := &models.Banner{Title: req.Title, TargetURL: req.TargetURL, SortOrder: req.SortOrder, IsActive: true}
	if req.IsActive != nil {
		b.IsActive = *req.IsActive
	}
	if err := h.bannerStore.Create(b); err != nil {
		utils.ServerError(w, h.logger, "create banner", err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "banner created successfully", "banner_id": b.ID})
}

// HandleUpdateBanner godoc
// @Summary      Update a banner (admin)
// @Tags         banners
// @Accept       json
// @Produce      json
// @Param        id   path string true "Banner ID"
// @Param        body body models.UpsertBannerRequest true "Banner"
// @Success      200 {object} MessageResponse
// @Security     BearerAuth
// @Router       /banners/update/{id} [put]
func (h *BannerHandler) HandleUpdateBanner(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	var req models.UpsertBannerRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	b := &models.Banner{ID: id, Title: req.Title, TargetURL: req.TargetURL, SortOrder: req.SortOrder, IsActive: true}
	if req.IsActive != nil {
		b.IsActive = *req.IsActive
	}
	if err = h.bannerStore.Update(b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "banner not found"})
			return
		}
		utils.ServerError(w, h.logger, "update banner", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "banner updated successfully"})
}

// HandleUpdateBannerImage godoc
// @Summary      Presigned URL for a banner image (admin)
// @Tags         banners
// @Produce      json
// @Param        id path string true "Banner ID"
// @Success      200 {object} PresignedURLResponse
// @Security     BearerAuth
// @Router       /banners/update/image/{id} [put]
func (h *BannerHandler) HandleUpdateBannerImage(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	key := fmt.Sprintf("banners/%s_%d.png", id, time.Now().Unix())
	url, err := h.blob.GenerateUploadPresignedURL(key)
	if err != nil {
		utils.ServerError(w, h.logger, "banner presign", err)
		return
	}
	if err = h.bannerStore.UpdateImage(id, "/"+key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "banner not found"})
			return
		}
		utils.ServerError(w, h.logger, "update banner image", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"url": url, "message": "banner presigned url generated successfully"})
}

// HandleDeleteBanner godoc
// @Summary      Delete a banner (admin)
// @Tags         banners
// @Produce      json
// @Param        id path string true "Banner ID"
// @Success      200 {object} MessageResponse
// @Security     BearerAuth
// @Router       /banners/delete/{id} [delete]
func (h *BannerHandler) HandleDeleteBanner(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	if err = h.bannerStore.Delete(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "banner not found"})
			return
		}
		utils.ServerError(w, h.logger, "delete banner", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "banner deleted successfully"})
}
