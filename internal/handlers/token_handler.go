package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

// TokenResponse represents a successful token response.
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type TokenHandler struct {
	userStore     store.UserStore
	businessStore store.BusinessStore
	logger        *slog.Logger
}

func NewTokenHandler(userStore store.UserStore, businessStore store.BusinessStore, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{
		userStore:     userStore,
		businessStore: businessStore,
		logger:        logger,
	}
}

// fullName returns a display name, handling a nil LastName safely.
func fullName(first string, last *string) string {
	if last == nil || *last == "" {
		return first
	}
	return fmt.Sprintf("%s %s", first, *last)
}

// HandleGetAdminTokenByEmailPassword godoc
// @Summary      Admin login
// @Description  Authenticates an admin using email and password and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "Admin login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid request payload"
// @Failure      401 {object} ErrorResponse "Invalid credentials"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /admin/login [post]
func (th *TokenHandler) HandleGetAdminTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.logger.Error("admin login", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	admin, err := th.userStore.GetAdminByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
			return
		}
		th.logger.Error("admin login", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	matched, err := admin.Password.Matches(req.Password)
	if err != nil {
		th.logger.Error("admin login: password comparison", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !matched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := tokens.GenerateNewToken(admin.ID, fullName(admin.FirstName, admin.LastName), nil)
	if err != nil {
		th.logger.Error("admin login: token generation", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}

// HandleGetUserTokenByEmailPassword godoc
// @Summary      User login
// @Description  Authenticates a user using email and password and returns a JWT token with business context
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body models.LoginRequest true "User login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid request payload"
// @Failure      401 {object} ErrorResponse "Invalid credentials"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /user/login [post]
func (th *TokenHandler) HandleGetUserTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.logger.Error("user login", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user, err := th.userStore.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
			return
		}
		th.logger.Error("user login", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	matched, err := user.Password.Matches(req.Password)
	if err != nil {
		th.logger.Error("user login: password comparison", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !matched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	businessID, err := th.businessStore.GetBusinessIDByUserID(user.ID)
	if err != nil {
		th.logger.Error("user login: fetch business id", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	token, err := tokens.GenerateNewToken(user.ID, fullName(user.FirstName, user.LastName), businessID)
	if err != nil {
		th.logger.Error("user login: token generation", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}
