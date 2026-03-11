package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// getTokenByEmailPasswordRequest is the login payload
type getTokenByEmailPasswordRequest struct {
	Email    string `json:"email"    example:"john@example.com"`
	Password string `json:"password" example:"StrongPass@1"`
}

// TokenResponse represents a successful token response
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

// HandleGetAdminTokenByEmailPassword godoc
// @Summary      Admin login
// @Description  Authenticates an admin using email and password, returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body getTokenByEmailPasswordRequest true "Admin login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or incorrect password"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /admin/login [post]
func (th *TokenHandler) HandleGetAdminTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req getTokenByEmailPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("ERROR: get admin token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	admin, err := th.userStore.GetAdminByEmail(req.Email)
	if err != nil {
		th.logger.Printf("ERROR get admin token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	isMatched, err := admin.Password.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR get admin token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !isMatched {
		th.logger.Println("ERROR get admin token by email password: incorrect password")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "incorrect password"})
		return
	}

	token, err := tokens.GenerateNewToken(admin.ID, fmt.Sprintf("%s %s", admin.FirstName, *admin.LastName), nil)
	if err != nil {
		th.logger.Printf("ERROR get admin token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}

// HandleGetUserTokenByEmailPassword godoc
// @Summary      User login
// @Description  Authenticates a user using email and password, returns a JWT token along with business context
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body getTokenByEmailPasswordRequest true "User login credentials"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or incorrect password"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /user/login [post]
func (th *TokenHandler) HandleGetUserTokenByEmailPassword(w http.ResponseWriter, r *http.Request) {
	var req getTokenByEmailPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("ERROR: get user token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := th.userStore.GetUserByEmail(req.Email)
	if err != nil {
		th.logger.Printf("ERROR get user token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	isMatched, err := user.Password.Matches(req.Password)
	if err != nil {
		th.logger.Printf("ERROR get user token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !isMatched {
		th.logger.Println("ERROR get user token by email password: incorrect password")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "incorrect password"})
		return
	}

	businessId, err := th.businessStore.GetBusinessIDByUserID(user.ID)
	if err != nil {
		th.logger.Printf("ERROR get user token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	token, err := tokens.GenerateNewToken(user.ID, fmt.Sprintf("%s %s", user.FirstName, *user.LastName), businessId)
	if err != nil {
		th.logger.Printf("ERROR get user token by email password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}