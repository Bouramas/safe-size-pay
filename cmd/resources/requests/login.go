package requests

import (
	"errors"
	"net/mail"
	"strings"
)

// LoginRequest represents the request made in order to login
// i.e. acquire a JWT token
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate performs basic validation checks on the LoginRequest struct.
func (r *LoginRequest) Validate() error {

	var required []string

	if r.Email == "" {
		required = append(required, "email")
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return errors.New("invalid email format")

	}

	if r.Password == "" {
		required = append(required, "password")
	}

	if len(required) > 0 {
		return errors.New(strings.Join(required, ", ") + " required.")
	}
	return nil
}
