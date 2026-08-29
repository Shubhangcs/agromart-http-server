package middlewares

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shubhangcs/agromart-server/internal/tokens"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

func claimsFrom(r *http.Request) *tokens.Token {
	c, _ := r.Context().Value("claims").(*tokens.Token)
	return c
}

// AdminOnly allows only tokens issued by /admin/login. Must run after AuthorizationMiddleware.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !claimsFrom(r).IsAdmin() {
			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "admin access required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SelfOrAdmin allows the request when the {id} path param is the caller's own user id, or the caller is an admin.
func SelfOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := claimsFrom(r)
		if c == nil || (!c.IsAdmin() && chi.URLParam(r, "id") != c.UserID) {
			utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "you can only modify your own account"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
