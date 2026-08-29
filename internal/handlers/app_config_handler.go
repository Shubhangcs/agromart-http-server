package handlers

import (
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/env"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// HandleAppConfig godoc
// @Summary      Public app configuration (minimum supported app version, store links)
// @Description  The app calls this on launch; if its version is below min_app_version it must prompt the user to update.
// @Tags         app
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /app/config [get]
func HandleAppConfig(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"min_app_version":    env.GetString("APP_MIN_VERSION", "1.0.0"),
		"latest_app_version": env.GetString("APP_LATEST_VERSION", "1.0.0"),
		"android_store_url":  env.GetString("APP_ANDROID_STORE_URL", "https://play.google.com/store/apps/details?id=com.southcanaraagromart.app"),
		"ios_store_url":      env.GetString("APP_IOS_STORE_URL", ""),
		"update_message":     env.GetString("APP_UPDATE_MESSAGE", "A new version of South Canara Agro Mart is available. Please update to continue."),
	})
}
