package handlers

import (
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/utils"
)

// businessLookup is the slice of BusinessStore needed to resolve a caller's business.
type businessLookup interface {
	GetBusinessIDByUserID(id string) (*string, error)
}

// callerOwnsBusiness reports whether the authenticated caller may act for businessID.
// Admins always may. The token's business_id claim is checked first; because a token
// issued before the user became a seller carries no business_id, the DB is consulted as
// a fallback so freshly created sellers are not locked out until they sign in again.
func callerOwnsBusiness(r *http.Request, bs businessLookup, businessID string) bool {
	c := claimsFromCtx(r)
	if c == nil || businessID == "" {
		return false
	}
	if c.IsAdmin() {
		return true
	}
	if c.BusinessID != nil && *c.BusinessID == businessID {
		return true
	}
	if bs != nil {
		if id, err := bs.GetBusinessIDByUserID(c.UserID); err == nil && id != nil && *id == businessID {
			return true
		}
	}
	return false
}

func forbidden(w http.ResponseWriter, msg string) {
	utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": msg})
}

const (
	msgNotYourBusiness = "you can only manage your own business"
	msgNotYourProduct  = "you can only manage your own products"
	msgNotYourRFQ      = "you can only manage your own RFQs"
)
