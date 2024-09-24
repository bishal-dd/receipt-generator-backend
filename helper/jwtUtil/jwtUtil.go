package jwtUtil

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokenSecret() string {
    secret := os.Getenv("SESSION_SECRET")
    if secret == "" {
        panic("SESSION_SECRET environment variable is not set")
    }
    return secret
}

// TokenFromRequest extracts the token from the Authorization header in the request
func TokenFromRequest(req *http.Request) (string, error) {
    authHeader := req.Header.Get("Authorization")
    if authHeader == "" {
        return "", errors.New("authorization header not found")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", errors.New("invalid authorization header format")
    }

    return parts[1], nil
}

// CreateToken generates a new JWT token for the given user
func CreateToken(user struct {
    ID        string
    Email     string
    FirstName string
    LastName  string
}) (string, error) {
    claims := jwt.MapClaims{
        "id":        user.ID,
        "email":     user.Email,
        "firstName": user.FirstName,
        "lastName":  user.LastName,
        "exp":       time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(GetTokenSecret()))
}

// ParseToken verifies and parses the token using the secret
func ParseToken(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(GetTokenSecret()), nil
    })
}

// DecodeToken decodes a JWT token without verifying it
func DecodeToken(tokenString string) (*jwt.Token, error) {
    return jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, nil)
}

func ClaimsFromToken(tokenString *jwt.Token) (jwt.MapClaims, error) {
    claims, ok := tokenString.Claims.(jwt.MapClaims)
		if !ok || !tokenString.Valid {
            return nil, errors.New("invalid token")
		}

    return claims, nil
}