package resources

import "github.com/golang-jwt/jwt/v5"

// Claims represents the JWT claims structure (payload).
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
