package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// getTokenByEmailPasswordRequest is the login payload.
type getTokenByEmailPasswordRequest struct {
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"StrongPass@1"`
}

// TokenResponse represents a successful token response.
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type TokenHandler struct {
	userStore     store.UserStore
	businessStore store.BusinessStore
	logger        *log.Logger
}

func NewTokenHandler(userStore store.UserStore, businessStore store.BusinessStore, logger *log.Logger) *TokenHandler {
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
// @Param        body body getTokenByEmailPasswordRequest true "Admin login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid request payload"
// @Failure      401 {object} ErrorResponse "Invalid credentials"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /admin/login [post]
func (th *TokenHandler) HandleGetAdminTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req getTokenByEmailPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.logger.Printf("ERROR: admin login: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	admin, err := th.userStore.GetAdminByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
			return
		}
		th.logger.Printf("ERROR: admin login: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	matched, err := admin.Password.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: admin login: password comparison: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !matched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := tokens.GenerateNewToken(admin.ID, fullName(admin.FirstName, admin.LastName), nil)
	if err != nil {
		th.logger.Printf("ERROR: admin login: token generation: %v\n", err)
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
// @Param        body body getTokenByEmailPasswordRequest true "User login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid request payload"
// @Failure      401 {object} ErrorResponse "Invalid credentials"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /user/login [post]
func (th *TokenHandler) HandleGetUserTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req getTokenByEmailPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		th.logger.Printf("ERROR: user login: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := th.userStore.GetUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
			return
		}
		th.logger.Printf("ERROR: user login: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	matched, err := user.Password.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: user login: password comparison: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if !matched {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	businessID, err := th.businessStore.GetBusinessIDByUserID(user.ID)
	if err != nil {
		th.logger.Printf("ERROR: user login: fetch business id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	token, err := tokens.GenerateNewToken(user.ID, fullName(user.FirstName, user.LastName), businessID)
	if err != nil {
		th.logger.Printf("ERROR: user login: token generation: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}
