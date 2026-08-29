// Package googleauth verifies Google ID tokens (the credential the mobile app receives from
// Google Sign-In) against Google's published signing keys, without a Google SDK dependency.
package googleauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shubhangcs/agromart-server/internal/env"
)

const certsURL = "https://www.googleapis.com/oauth2/v3/certs"

// Identity is what we trust from a verified token.
type Identity struct {
	Sub, Email, GivenName, FamilyName, Picture string
	EmailVerified                              bool
}

type Verifier interface {
	Verify(ctx context.Context, idToken string) (*Identity, error)
}

// New returns a JWKS-backed verifier for the configured client IDs (GOOGLE_OAUTH_CLIENT_IDS,
// comma-separated). GOOGLE_AUTH_TEST_MODE=true swaps in a fake that accepts
// "test:<sub>:<email>:<given>:<family>" — for automated tests only, never in production.
func New() Verifier {
	if strings.EqualFold(env.GetString("GOOGLE_AUTH_TEST_MODE", "false"), "true") {
		return fakeVerifier{}
	}
	var auds []string
	for _, a := range strings.Split(env.GetString("GOOGLE_OAUTH_CLIENT_IDS", ""), ",") {
		if a = strings.TrimSpace(a); a != "" {
			auds = append(auds, a)
		}
	}
	return &googleVerifier{audiences: auds, client: &http.Client{Timeout: 10 * time.Second}}
}

type claims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	jwt.RegisteredClaims
}

type googleVerifier struct {
	audiences []string
	client    *http.Client

	mu      sync.Mutex
	keys    map[string]*rsa.PublicKey
	expires time.Time
}

func (v *googleVerifier) Verify(ctx context.Context, idToken string) (*Identity, error) {
	if len(v.audiences) == 0 {
		return nil, errors.New("google sign-in is not configured on the server")
	}
	var c claims
	tok, err := jwt.ParseWithClaims(idToken, &c, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		kid, _ := t.Header["kid"].(string)
		return v.key(ctx, kid)
	}, jwt.WithIssuer("https://accounts.google.com"), jwt.WithLeeway(30*time.Second))
	if err != nil || !tok.Valid {
		// Google also issues tokens with iss = "accounts.google.com"
		tok, err = jwt.ParseWithClaims(idToken, &c, func(t *jwt.Token) (any, error) {
			kid, _ := t.Header["kid"].(string)
			return v.key(ctx, kid)
		}, jwt.WithIssuer("accounts.google.com"), jwt.WithLeeway(30*time.Second))
		if err != nil || !tok.Valid {
			return nil, fmt.Errorf("invalid google token: %w", err)
		}
	}
	audOK := false
	for _, want := range v.audiences {
		for _, got := range c.Audience {
			if got == want {
				audOK = true
			}
		}
	}
	if !audOK {
		return nil, errors.New("google token was issued for a different app")
	}
	if c.Subject == "" || c.Email == "" {
		return nil, errors.New("google token missing subject or email")
	}
	return &Identity{Sub: c.Subject, Email: strings.ToLower(c.Email), EmailVerified: c.EmailVerified,
		GivenName: c.GivenName, FamilyName: c.FamilyName, Picture: c.Picture}, nil
}

func (v *googleVerifier) key(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if k, ok := v.keys[kid]; ok && time.Now().Before(v.expires) {
		return k, nil
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, certsURL, nil)
	res, err := v.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var jwks struct {
		Keys []struct{ Kid, N, E string }
	}
	if err = json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		return nil, err
	}
	keys := map[string]*rsa.PublicKey{}
	for _, k := range jwks.Keys {
		n, err1 := base64.RawURLEncoding.DecodeString(k.N)
		e, err2 := base64.RawURLEncoding.DecodeString(k.E)
		if err1 != nil || err2 != nil {
			continue
		}
		keys[k.Kid] = &rsa.PublicKey{N: new(big.Int).SetBytes(n), E: int(new(big.Int).SetBytes(e).Int64())}
	}
	v.keys, v.expires = keys, time.Now().Add(6*time.Hour)
	k, ok := keys[kid]
	if !ok {
		return nil, errors.New("unknown google signing key")
	}
	return k, nil
}

type fakeVerifier struct{}

func (fakeVerifier) Verify(_ context.Context, idToken string) (*Identity, error) {
	p := strings.Split(idToken, ":")
	if len(p) < 3 || p[0] != "test" {
		return nil, errors.New("invalid google token")
	}
	id := &Identity{Sub: p[1], Email: strings.ToLower(p[2]), EmailVerified: true}
	if len(p) > 3 {
		id.GivenName = p[3]
	}
	if len(p) > 4 {
		id.FamilyName = p[4]
	}
	return id, nil
}
