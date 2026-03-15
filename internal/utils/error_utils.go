package utils

import (
	"log/slog"
	"net/http"
)

func ServerError(w http.ResponseWriter, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	WriteJSON(w, http.StatusInternalServerError, Envelope{"error": "internal server error"})
}

func BadRequest(w http.ResponseWriter, logger *slog.Logger, msg string, err error) {
	logger.Error(msg, "error", err)
	WriteJSON(w, http.StatusBadRequest, Envelope{"error": msg})
}
