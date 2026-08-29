package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/shubhangcs/agromart-server/internal/blob"
	"github.com/shubhangcs/agromart-server/internal/googleauth"
	"github.com/shubhangcs/agromart-server/internal/mailer"
	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

type SocialAuthHandler struct {
	userStore     store.UserStore
	businessStore store.BusinessStore
	blobStore     store.BlobStore
	blob          *blob.AWSS3
	verifier      googleauth.Verifier
	mailer        mailer.Mailer
	logger        *slog.Logger
}

func NewSocialAuthHandler(userStore store.UserStore, businessStore store.BusinessStore, blobStore store.BlobStore, b *blob.AWSS3, v googleauth.Verifier, m mailer.Mailer, logger *slog.Logger) *SocialAuthHandler {
	return &SocialAuthHandler{userStore: userStore, businessStore: businessStore, blobStore: blobStore, blob: b, verifier: v, mailer: m, logger: logger}
}

type googleAuthRequest struct {
	IDToken string `json:"id_token"`
}

// HandleGoogleAuth godoc
// @Summary      Sign in or sign up with Google
// @Description  Verifies a Google ID token. Links to an existing account with the same verified email, or creates a new password-less account. Returns the app JWT plus needs_phone (true until the user adds a phone number).
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body googleAuthRequest true "Google ID token"
// @Success      200 {object} map[string]interface{}
// @Failure      401 {object} ErrorResponse
// @Router       /user/auth/google [post]
func (h *SocialAuthHandler) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req googleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IDToken == "" {
		utils.BadRequest(w, h.logger, "id_token is required", err)
		return
	}
	id, err := h.verifier.Verify(r.Context(), req.IDToken)
	if err != nil {
		h.logger.Warn("google auth", "error", err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "google sign-in could not be verified"})
		return
	}
	if !id.EmailVerified {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "your google email address is not verified"})
		return
	}

	isNew := false
	user, err := h.userStore.GetUserByGoogleSub(id.Sub)
	if errors.Is(err, sql.ErrNoRows) {
		// Not linked yet: attach to an existing account with the same (verified) email, or create one.
		user, err = h.userStore.GetUserByEmail(id.Email)
		switch {
		case err == nil:
			if err = h.userStore.LinkGoogle(user.ID, id.Sub); err != nil {
				utils.ServerError(w, h.logger, "link google account", err)
				return
			}
		case errors.Is(err, sql.ErrNoRows):
			first := id.GivenName
			if first == "" {
				first = "Trader"
			}
			user = &models.User{FirstName: first, Email: id.Email, AuthProvider: "google"}
			if id.FamilyName != "" {
				ln := id.FamilyName
				user.LastName = &ln
			}
			if id.Picture != "" {
				pic := id.Picture
				user.ProfileImage = &pic
			}
			sub := id.Sub
			user.GoogleSub = &sub
			if err = h.userStore.CreateGoogleUser(user); err != nil {
				utils.ServerError(w, h.logger, "create google user", err)
				return
			}
			isNew = true
			mailer.SendWelcome(h.mailer, h.logger, user.Email, user.FirstName)
			if id.Picture != "" {
				go h.copyAvatarToS3(user.ID, id.Picture)
			}
		default:
			utils.ServerError(w, h.logger, "google auth lookup", err)
			return
		}
	} else if err != nil {
		utils.ServerError(w, h.logger, "google auth lookup", err)
		return
	}
	if user.IsUserBlocked {
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "this account has been blocked"})
		return
	}

	businessID, _ := h.businessStore.GetBusinessIDByUserID(user.ID)
	name := user.FirstName
	if user.LastName != nil && *user.LastName != "" {
		name += " " + *user.LastName
	}
	token, err := tokens.GenerateNewToken(user.ID, name, tokens.RoleUser, businessID)
	if err != nil {
		utils.ServerError(w, h.logger, "generate token", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"token":       token,
		"user_id":     user.ID,
		"is_new_user": isNew,
		"needs_phone": user.Phone == "",
	})
}

// copyAvatarToS3 mirrors the Google profile picture into our bucket so the app never depends on
// Google's CDN (their URLs are size-limited and can expire). Failures keep the Google URL.
func (h *SocialAuthHandler) copyAvatarToS3(userID, pictureURL string) {
	if h.blob == nil || h.blobStore == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// ask Google for a larger rendition than the default 96px
	src := strings.Replace(pictureURL, "=s96-c", "=s512-c", 1)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, src, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		h.logger.Warn("google avatar fetch failed", "user_id", userID, "error", err)
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(io.LimitReader(res.Body, 5<<20))
	if err != nil || len(data) == 0 {
		return
	}
	ct := res.Header.Get("Content-Type")
	if ct == "" {
		ct = "image/jpeg"
	}
	key := fmt.Sprintf("profile/user/%s_%d.png", userID, time.Now().Unix())
	if err = h.blob.UploadObject(ctx, key, ct, data); err != nil {
		h.logger.Warn("google avatar upload failed", "user_id", userID, "error", err)
		return
	}
	if err = h.blobStore.UpdateUserProfileImage(userID, "/"+key); err != nil {
		h.logger.Warn("google avatar db update failed", "user_id", userID, "error", err)
	}
}
