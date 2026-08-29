package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/shubhangcs/agromart-server/internal/mailer"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

const (
	resetCodeTTL         = 15 * time.Minute
	resetCodeMaxAttempts = 5
)

type PasswordResetHandler struct {
	userStore  store.UserStore
	resetStore store.PasswordResetStore
	mailer     mailer.Mailer
	logger     *slog.Logger
}

func NewPasswordResetHandler(userStore store.UserStore, resetStore store.PasswordResetStore, m mailer.Mailer, logger *slog.Logger) *PasswordResetHandler {
	return &PasswordResetHandler{userStore: userStore, resetStore: resetStore, mailer: m, logger: logger}
}

type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type resetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6,numeric"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// HandleForgotPassword godoc
// @Summary      Request a password-reset code
// @Description  Emails a 6-digit code (valid 15 minutes) if an account exists. Always returns 200 so email addresses cannot be enumerated.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body forgotPasswordRequest true "Email"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Router       /user/forgot-password [post]
func (h *PasswordResetHandler) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	generic := utils.Envelope{"message": "if an account exists for this email, a reset code has been sent"}

	user, err := h.userStore.GetUserByEmail(email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("forgot password lookup", "error", err)
		}
		utils.WriteJSON(w, http.StatusOK, generic)
		return
	}

	if len(user.Password.Hash) == 0 {
		// Google-only account: nothing to reset. Still 200 so emails cannot be enumerated.
		utils.WriteJSON(w, http.StatusOK, generic)
		return
	}
	code, err := randomCode()
	if err != nil {
		utils.ServerError(w, h.logger, "generate reset code", err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(code), 10)
	if err != nil {
		utils.ServerError(w, h.logger, "hash reset code", err)
		return
	}
	if err = h.resetStore.Create(email, string(hash), time.Now().Add(resetCodeTTL)); err != nil {
		utils.ServerError(w, h.logger, "store reset code", err)
		return
	}

	subject, text, html := mailer.ResetCodeEmail(user.FirstName, code)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if err = h.mailer.Send(ctx, email, subject, text, html); err != nil {
		utils.ServerError(w, h.logger, "send reset email", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, generic)
}

// HandleResetPassword godoc
// @Summary      Reset password with an emailed code
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body resetPasswordRequest true "Email, code, new password"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Router       /user/reset-password [post]
func (h *PasswordResetHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, h.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, h.logger, err.Error(), err)
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))
	invalid := func() {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired code"})
	}

	pr, err := h.resetStore.GetActive(email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("reset lookup", "error", err)
		}
		invalid()
		return
	}
	if pr.Attempts >= resetCodeMaxAttempts {
		_ = h.resetStore.MarkUsed(pr.ID)
		invalid()
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(pr.CodeHash), []byte(req.Code)) != nil {
		_ = h.resetStore.IncrementAttempts(pr.ID)
		invalid()
		return
	}

	user, err := h.userStore.GetUserByEmail(email)
	if err != nil {
		invalid()
		return
	}
	if err = user.Password.Set(req.NewPassword); err != nil {
		utils.ServerError(w, h.logger, "hash new password", err)
		return
	}
	if err = h.userStore.UpdateUserPassword(user); err != nil {
		utils.ServerError(w, h.logger, "update password", err)
		return
	}
	_ = h.resetStore.MarkUsed(pr.ID)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "password reset successfully"})
}

func randomCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
