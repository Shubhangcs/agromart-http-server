package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/push"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

type PushHandler struct {
	pushStore store.PushStore
	logger    *slog.Logger
}

func NewPushHandler(pushStore store.PushStore, logger *slog.Logger) *PushHandler {
	return &PushHandler{pushStore: pushStore, logger: logger}
}

type pushTokenRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

// HandleRegisterPushToken godoc
// @Summary      Register this device for push notifications
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body pushTokenRequest true "Expo push token"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /user/push-token [post]
func (h *PushHandler) HandleRegisterPushToken(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r)
	var req pushTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !push.IsExpoToken(req.Token) {
		utils.BadRequest(w, h.logger, "a valid Expo push token is required", err)
		return
	}
	platform := req.Platform
	if platform != "ios" && platform != "android" {
		platform = "android"
	}
	if err := h.pushStore.Upsert(claims.UserID, req.Token, platform); err != nil {
		utils.ServerError(w, h.logger, "register push token", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "push token registered"})
}

// HandleUnregisterPushToken godoc
// @Summary      Remove this device's push token (call on logout)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body pushTokenRequest true "Expo push token"
// @Success      200 {object} MessageResponse
// @Security     BearerAuth
// @Router       /user/push-token [delete]
func (h *PushHandler) HandleUnregisterPushToken(w http.ResponseWriter, r *http.Request) {
	var req pushTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Token == "" {
		utils.BadRequest(w, h.logger, "token is required", err)
		return
	}
	if err := h.pushStore.Delete(req.Token); err != nil {
		utils.ServerError(w, h.logger, "unregister push token", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "push token removed"})
}
