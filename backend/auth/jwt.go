package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWTSECRET")
	if secret == "" {
		panic("JWTSECRET not set in environment")
	}
	return []byte(secret)
}

// generate jwt for signed token for a user
func GenerateJWT(userID uint32, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// validats signed jwt
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// parses jwt
	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (any, error) {
			// ensure algorithm is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return getJWTSecret(), nil
		})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return token, nil
}

// extracts user id and role from a validated token
func ExtractClaims(token *jwt.Token) (uint32, string, error) {
	// get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("could not parse claims")
	}

	// safe assertions
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("user id not found in token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("role not found in token")
	}

	return uint32(userIDFloat), role, nil
}
