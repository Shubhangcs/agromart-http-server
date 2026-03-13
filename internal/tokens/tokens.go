package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	ttl          = time.Hour * 24 // 24-hour access token
	issuer       = os.Getenv("JWT_TOKEN_ISSUER")
)

// Token holds the JWT claims used throughout the application.
type Token struct {
	UserID     string  `json:"user_id"`
	UserName   string  `json:"user_name"`
	BusinessID *string `json:"business_id,omitempty"`
	jwt.RegisteredClaims
}

// GenerateNewToken signs and returns a new JWT for the given user.
func GenerateNewToken(userID, userName string, businessID *string) (string, error) {
	claims := Token{
		UserID:     userID,
		UserName:   userName,
		BusinessID: businessID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

// ValidateToken parses and validates the given JWT string.
func ValidateToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Token{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(jwtSecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Token)
	if !ok || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	if claims.Issuer != issuer {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}
