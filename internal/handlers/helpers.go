package handlers

import (
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/utils"
)

// serverError logs the error and writes a 500 response.
func serverError(w http.ResponseWriter, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
}

// badRequest writes a 400 response with the given message.
func badRequest(w http.ResponseWriter, msg string) {
	utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": msg})
}
