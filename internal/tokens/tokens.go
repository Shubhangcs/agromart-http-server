package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	ttl          = time.Hour * 24 * 365
	issuer       = os.Getenv("JWT_TOKEN_ISSUER")
)

type Token struct {
	UserID     string  `json:"user_id"`
	UserName   string  `json:"user_name"`
	BusinessID *string `json:"business_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateNewToken(userId string, userName string, businessId *string) (string, error) {
	claims := Token{
		UserID:     userId,
		UserName:   userName,
		BusinessID: businessId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Token{},
		func(t *jwt.Token) (any, error) {
			return jwtSecretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Token)
	if !ok {
		return nil, errors.New("invalid authorization token or token is expired")
	}

	if claims.Issuer != issuer {
		return nil, errors.New("invalid authorization token")
	}

	return claims, nil
}
