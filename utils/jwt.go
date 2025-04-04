package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Use environment variables for better security
var jwtSecret = []byte(getEnv("JWT_SECRET", "nnovrian"))

// Claims struct for JWT
type Claims struct {
	UserID      uint   `json:"user_id"`
	Role        string `json:"role"`
	CompanyName string `json:"company_name"`
	OwnerID     *uint  `json:"owner_id"` // New attribute, optional
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token
func GenerateToken(userID uint, role string, companyName string, duration time.Duration) (string, error) {
	claims := Claims{
		UserID:      userID,
		Role:        role,
		CompanyName: companyName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken verifies and decodes a JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		var vErr *jwt.ValidationError
		if errors.As(err, &vErr) {
			switch {
			case vErr.Errors&jwt.ValidationErrorExpired != 0:
				return nil, errors.New("token expired")
			case vErr.Errors&jwt.ValidationErrorSignatureInvalid != 0:
				return nil, errors.New("invalid token signature")
			default:
				return nil, errors.New("token is invalid")
			}
		}
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// getEnv retrieves environment variables or returns a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
